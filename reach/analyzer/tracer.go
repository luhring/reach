package analyzer

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/luhring/reach/reach"
)

type Tracer struct {
	provider reach.InfrastructureGetter
}

func NewTracer(provider reach.InfrastructureGetter) *Tracer {
	return &Tracer{
		provider: provider,
	}
}

func (t *Tracer) Trace(source, destination reach.Subject) ([]reach.Path, error) {
	sourceRef := reach.InfrastructureReference{
		R: reach.ResourceReference{
			Domain: source.Domain,
			Kind:   source.Kind,
			ID:     source.ID,
		},
	}

	destinationRef := reach.InfrastructureReference{
		R: reach.ResourceReference{
			Domain: destination.Domain,
			Kind:   destination.Kind,
			ID:     destination.ID,
		},
	}

	initialJob := traceJob{
		ref:            sourceRef,
		destinationRef: destinationRef,
	}

	done := make(chan interface{})
	defer close(done)

	results := t.tracePoint(done, initialJob)
	var paths []reach.Path

	for result := range results {
		if result.error != nil {
			_, _ = fmt.Fprintln(os.Stderr, result.error) // TODO: Log more intelligently!
		}
		paths = append(paths, *result.path)
	}

	return paths, nil
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
				resource, err := t.provider.Get(job.ref)
				if err != nil {
					results <- traceResult{error: err}
					return
				}
				traceable := resource.Properties.(reach.Traceable) // TODO: Make this more type-safe

				err = detectLoop(job.path, traceable)
				if err != nil {
					results <- traceResult{error: fmt.Errorf("tracer detected a loop: %v", err)}
					return
				}

				factors := traceable.Factors()
				point := reach.Point{
					Ref:     job.ref,
					Factors: factors,
				}

				var path reach.Path
				var destIPs []net.IP
				firstTrace := job.path == nil
				if firstTrace {
					// This is the first traced point.
					path = reach.NewPath(point)
				} else {
					path = *job.path
					path.Add(job.edgeTuple, point, traceable.Segments())
				}

				if traceable.Ref().Equal(job.destinationRef) {
					// Path is complete!
					results <- traceResult{path: &path}
					return
				}

				// TODO: Maybe this is where we should account for multiple dstIPs

				var edgeTuples []reach.IPTuple

				if firstTrace {
					destResource, err := t.provider.Get(job.destinationRef)
					if err != nil {
						results <- traceResult{error: fmt.Errorf("unable to get destination: %v", err)}
						return
					}
					dest := destResource.Properties.(reach.IPAdvertiser)
					destIPs, err = dest.IPs(t.provider)
					if err != nil {
						results <- traceResult{error: fmt.Errorf("unable to get advertised IPs from destination: %v", err)}
						return
					}
					for _, dst := range destIPs {
						edgeTuples = append(edgeTuples, reach.IPTuple{
							Src: nil, // TODO: There are plenty of cases where this isn't nil (such as generic domain IP address as source)
							Dst: dst,
						})
					}
				} else {
					if job.edgeTuple != nil {
						edgeTuples = []reach.IPTuple{*job.edgeTuple}
					}
				}

				var edges []reach.PathEdge
				for _, tuple := range edgeTuples {
					tupleEdges, err := traceable.ForwardEdges(&tuple, t.provider)
					if err != nil {
						results <- traceResult{error: fmt.Errorf("tracer was unable to get edges for ref (%s): %v", job.ref, err)}
						return
					}
					edges = append(edges, tupleEdges...)
				}

				numEdges := len(edges)
				if numEdges < 1 {
					err := fmt.Errorf("no forward edges found when processing job:\n%+v", job)
					results <- traceResult{error: err}
					return
				}

				resultChannels := make([]<-chan traceResult, numEdges)
				for _, edge := range edges {
					j := traceJob{
						ref:            edge.Ref,
						path:           &path,
						edgeTuple:      edge.Tuple,
						destinationRef: job.destinationRef,
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

func detectLoop(path *reach.Path, traceable reach.Traceable) error {
	// TODO: (Later) Consider a more intelligent loop detection system that leverages tuples

	if path == nil {
		return nil
	}

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
