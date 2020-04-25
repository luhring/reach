package analyzer

import (
	"fmt"
	"net"
	"sync"

	"github.com/luhring/reach/reach"
)

type Tracer struct {
	referenceResolver    reach.ReferenceResolver
	domainClientResolver reach.DomainClientResolver
}

func NewTracer(domainClientResolver reach.DomainClientResolver) *Tracer {
	referenceResolver := NewReferenceResolver(domainClientResolver)

	return &Tracer{
		referenceResolver:    &referenceResolver,
		domainClientResolver: domainClientResolver,
	}
}

func (t *Tracer) Trace(source, destination reach.Subject) ([]reach.Path, error) {
	dstIPs, err := t.subjectIPs(destination)
	if err != nil {
		return nil, fmt.Errorf("trace failed: unable to get IPs for destination: %v", err)
	}

	initialJob := traceJob{
		ref:            source.Ref(),
		sourceRef:      source.Ref(),
		destinationRef: destination.Ref(),
		destinationIPs: dstIPs,
	}

	done := make(chan interface{})
	defer close(done)

	results := t.tracePoint(done, initialJob)
	var paths []reach.Path

	for result := range results {
		if result.error != nil {
			// _, _ = fmt.Fprintln(os.Stderr, result.error) // TODO: Log more intelligently!
			return nil, result.error
		}
		paths = append(paths, *result.path)
	}

	// TODO: Backtrace all paths to fill in return factors

	return paths, nil
}

func (t *Tracer) subjectIPs(s reach.Subject) ([]net.IP, error) {
	infrastructure, err := t.referenceResolver.Resolve(s.Ref())
	if err != nil {
		return nil, fmt.Errorf("unable to get infrastructure for subject (%v): %v", s.Ref(), err)
	}
	addressable, ok := infrastructure.Properties.(reach.IPAddressable)
	if !ok {
		return nil, fmt.Errorf("subject does not implement IPAddressable (%v)", s.Ref())
	}
	ips, err := addressable.IPs(t.domainClientResolver)
	if err != nil {
		return nil, fmt.Errorf("unable to get IP addresses for subject (%v): %v", s.Ref(), err)
	}
	return ips, nil
}

func (t *Tracer) tracePoint(done <-chan interface{}, job traceJob) <-chan traceResult {
	results := make(chan traceResult)

	go func() {
		defer close(results)
		for {
			select {
			case <-done:
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
					results <- traceResult{error: fmt.Errorf("tracer detected a path cycle: %v", err)}
					return
				}

				isFirstTraceJob := job.path.Zero()

				var previousEdge *reach.Edge
				if isFirstTraceJob == false {
					previousEdge = &job.edge
				}

				factors, err := traceable.FactorsForward(t.domainClientResolver, previousEdge)
				point := reach.Point{Ref: job.ref, Factors: factors}

				var path reach.Path
				if isFirstTraceJob {
					path = reach.NewPath(point)
				} else {
					path = job.path.Add(job.edge, point, traceable.Segments())
				}

				if traceable.Ref().Equal(job.destinationRef) {
					// Path is complete!
					results <- traceResult{path: &path}
					return
				}

				edges, err := traceable.EdgesForward(t.domainClientResolver, previousEdge, job.destinationIPs)
				if err != nil {
					results <- traceResult{
						error: fmt.Errorf("tracer was unable to get edges for ref (%s): %v", job.ref, err),
					}
					return
				}

				numEdges := len(edges)
				if numEdges < 1 {
					results <- traceResult{
						error: fmt.Errorf("no forward edges found when processing job:\n%+v", job)}
					return
				}

				resultChannels := make([]<-chan traceResult, numEdges)
				for _, edge := range edges {
					j := traceJob{
						ref:            edge.EndRef,
						path:           path,
						edge:           edge,
						sourceRef:      job.sourceRef,
						destinationRef: job.destinationRef,
						destinationIPs: job.destinationIPs,
					}
					resultChannels = append(resultChannels, t.tracePoint(done, j))
				}

				// Wait for downstream results to come in and pass them upstream.
				downstreamResults := fanIn(done, resultChannels)
				for r := range downstreamResults {
					select {
					case <-done:
						return
					case results <- r:
					}
				}
			}
		}
	}()

	return results
}

func ensureNoPathCycles(path reach.Path, traceable reach.Traceable) error {
	// TODO: (Later) Consider a more intelligent cycle detection system that leverages tuples

	ref := traceable.Ref()
	if traceable.Visitable(path.Contains(ref)) == false {
		return fmt.Errorf("cannot visit point again: %s", ref)
	}
	return nil
}

func fanIn(done <-chan interface{}, channels []<-chan traceResult) <-chan traceResult {
	var wg sync.WaitGroup
	multiplexedStream := make(chan traceResult)

	multiplex := func(c <-chan traceResult) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
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
