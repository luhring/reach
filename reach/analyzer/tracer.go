package analyzer

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/reacherr"
	"github.com/luhring/reach/reach/reachlog"
)

// Tracer is the analyzer-specific implementation of the interface reach.Tracer. This implementation features a mechanism for tracing paths that concurrently follows all possible paths of network traffic from source to destination.
type Tracer struct {
	referenceResolver    *ReferenceResolver
	domainClientResolver reach.DomainClientResolver
	logger               reachlog.Logger
}

// NewTracer returns a reference to a new instance of a Tracer.
func NewTracer(domainClientResolver reach.DomainClientResolver, logger reachlog.Logger) *Tracer {
	referenceResolver := NewReferenceResolver(domainClientResolver)

	return &Tracer{
		referenceResolver:    &referenceResolver,
		domainClientResolver: domainClientResolver,
		logger:               logger,
	}
}

// Trace uses available information to map all possible network paths from the specified source to the specified destination. If Trace is unable to provide a complete set of paths, it returns an error.
func (t *Tracer) Trace(source, destination reach.Subject) ([]reach.Path, error) {
	t.logger.Debug("beginning trace", "source", source, "destination", destination)

	dstIPs, err := t.subjectIPs(destination)
	if err != nil {
		t.logger.Error("trace failed: unable to get IPs for destination", "err", err, "destination", destination.Ref())
		return nil, err
	}

	initialJob := traceJob{
		ref:            source.Ref(),
		sourceRef:      source.Ref(),
		destinationRef: destination.Ref(),
		destinationIPs: dstIPs,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := t.tracePoint(ctx, initialJob)
	var paths []reach.Path

	for result := range results {
		if result.error != nil {
			t.logger.Error("trace failed: error while tracing point", "err", result.error)
			return nil, result.error
		}
		paths = append(paths, *result.path)
	}

	t.logger.Debug("beginning return factors calculation")

	for i, path := range paths {
		updatedPath, err := path.MapPoints(func(point reach.Point, leftEdge, rightEdge *reach.Edge) (reach.Point, error) {
			resource, err := t.referenceResolver.Resolve(point.Ref)
			if err != nil {
				return reach.Point{}, err
			}
			traceable, ok := resource.Properties.(reach.Traceable)
			if !ok {
				return reach.Point{}, fmt.Errorf("cannot return-trace point that doesn't implement the traceable interface: '%v'", point.Ref)
			}

			factorsReturn, err := traceable.FactorsReturn(t.domainClientResolver, rightEdge)
			if err != nil {
				return reach.Point{}, err
			}

			return reach.Point{
				Ref:            point.Ref,
				FactorsForward: point.FactorsForward,
				FactorsReturn:  factorsReturn,
				SegmentDivider: point.SegmentDivider,
			}, nil
		})

		if err != nil {
			return nil, err
		}
		paths[i] = updatedPath
	}

	t.logger.Info("trace successful", "numPaths", len(paths))
	return paths, nil
}

func (t *Tracer) subjectIPs(s reach.Subject) ([]net.IP, error) {
	subjectResource, err := t.referenceResolver.Resolve(s.Ref())
	if err != nil {
		return nil, err
	}
	addressable, ok := subjectResource.Properties.(reach.IPAddressable)
	if !ok {
		msg := "subject does not implement IPAddressable"
		t.logger.Error(msg, "subject", s.Ref())
		err = fmt.Errorf(msg+": %v", s.Ref())
		return nil, err
	}
	ips, err := addressable.IPs(t.domainClientResolver)
	if err != nil {
		t.logger.Error("unable to get subject IPs", "err", err, "subject", s.Ref())
		return nil, err
	}
	return ips, nil
}

func (t *Tracer) tracePoint(ctx context.Context, job traceJob) <-chan traceResult {
	t.logger.Info("tracing point", "ref", job.ref)

	results := make(chan traceResult)

	go func() {
		defer close(results)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// We need to turn the ref into a Traceable
				resource, err := t.referenceResolver.Resolve(job.ref)
				if err != nil {
					results <- traceResult{error: err}
					return
				}
				traceable, ok := resource.Properties.(reach.Traceable)
				if !ok {
					results <- traceResult{
						error: fmt.Errorf("cannot trace point that doesn't implement the traceable interface: '%v'", job.ref),
					}
					return
				}

				err = ensureNoPathCycles(job.path, traceable)
				if err != nil {
					t.logger.Error("path cycle detected", "err", err)
					results <- traceResult{error: err}
					return
				}

				isFirstTraceJob := job.path.Zero()

				var previousEdge *reach.Edge
				var previousRef *reach.Reference
				if isFirstTraceJob == false {
					previousEdge = &job.edge
					r := job.path.LastPoint().Ref
					previousRef = &r
				}

				factors, err := traceable.FactorsForward(t.domainClientResolver, previousEdge)
				t.logger.Debug("discovered factors", "ref", job.ref, "numFactors", len(factors))
				point := reach.Point{Ref: job.ref, FactorsForward: factors, SegmentDivider: traceable.Segments()}

				var path reach.Path
				if isFirstTraceJob {
					path = reach.NewPath(point)
				} else {
					path = job.path.Add(job.edge, point)
				}

				if traceable.Ref().Equal(job.destinationRef) {
					// Path is complete!
					t.logger.Info("completed trace of path", "source", job.sourceRef, "destination", job.destinationRef)

					results <- traceResult{path: &path}
					return
				}

				edges, err := traceable.EdgesForward(t.domainClientResolver, previousEdge, previousRef, job.destinationIPs)
				if err != nil {
					t.logger.Error("tracer was unable to get edges forward", "err", err, "ref", job.ref)
					results <- traceResult{
						error: err,
					}
					return
				}

				numEdges := len(edges)
				if numEdges < 1 {
					msg := "no forward edges found when processing job"
					err = reacherr.New(nil, msg+":\n%+v", job)
					t.logger.Error(msg, "job", job)
					results <- traceResult{
						error: err,
					}
					return
				}
				t.logger.Debug("found edge(s) forward from trace point", "ref", job.ref, "numEdges", numEdges)

				resultChannels := make([]<-chan traceResult, numEdges)
				for i, edge := range edges {
					nextJob := traceJob{
						ref:            edge.EndRef,
						path:           path,
						edge:           edge,
						sourceRef:      job.sourceRef,
						destinationRef: job.destinationRef,
						destinationIPs: job.destinationIPs,
					}

					resultChannels[i] = t.tracePoint(ctx, nextJob)
				}

				// Wait for downstream results to come in and pass them upstream.
				downstreamResults := fanIn(ctx, resultChannels)
				for r := range downstreamResults {
					select {
					case <-ctx.Done():
						return
					case results <- r:
					}
				}
				return
			}
		}
	}()

	return results
}

func ensureNoPathCycles(path reach.Path, traceable reach.Traceable) error {
	// TODO: (Later) Consider a more intelligent cycle detection system that leverages tuples

	ref := traceable.Ref()
	if traceable.Visitable(path.Contains(ref)) == false {
		return reacherr.New(nil, "cannot visit point again: %s", ref)
	}
	return nil
}

func fanIn(ctx context.Context, channels []<-chan traceResult) <-chan traceResult {
	var wg sync.WaitGroup
	multiplexedStream := make(chan traceResult)

	multiplex := func(c <-chan traceResult) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}
