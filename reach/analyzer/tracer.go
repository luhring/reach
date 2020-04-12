package analyzer

import (
	"fmt"
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

	initialJob := tracerJob{
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

func (t *Tracer) tracePoint(done <-chan interface{}, job tracerJob) <-chan traceResult {
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
				if job.path == nil {
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

				var wg sync.WaitGroup

				edges, err := traceable.ForwardEdges(job.edgeTuple, t.provider)
				if err != nil {
					results <- traceResult{error: fmt.Errorf("tracer was unable to get edges for ref (%s): %v", job.ref, err)}
					return
				}
				numEdges := len(edges)
				if numEdges < 1 {
					err := fmt.Errorf("no next points found when processing job:\n%+v", job)
					results <- traceResult{error: err}
					return
				}

				resultChannels := make([]<-chan traceResult, numEdges)
				wg.Add(numEdges)
				for _, edge := range edges {
					j := tracerJob{
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
