package analyzer

import (
	"fmt"
	"os"
	"sync"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/analyzer/pathbuilder"
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

	dest, err := t.provider.Get(destinationRef)
	if err != nil {
		return nil, fmt.Errorf("tracer could not get destination: %v", err)
	}
	destIPResolver := dest.Properties.(reach.SubjectIPResolver)
	destIPs, err := destIPResolver.Resolve(reach.SubjectRoleDestination, t.provider)
	if err != nil {
		return nil, fmt.Errorf("couldn't resolve IPs for destination (%s): %v", destination, err)
	}

	initialJob := tracerJob{
		partial: reach.PartialPath{
			Path:    reach.NewPath(),
			NextRef: sourceRef,
		},
		destinationIPs: destIPs,
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
				partial := job.partial

				// We need to turn the ref into a Traceable
				r, err := t.provider.Get(partial.NextRef)
				if err != nil {
					results <- traceResult{error: err}
					return
				}
				traceable := r.Properties.(reach.Traceable) // TODO: Make this more type-safe

				factors := traceable.Factors()
				prevTuple := partial.Path.LastPoint().Tuple
				newTuple := traceable.NextTuple(prevTuple)
				newPoint := reach.Point{
					Ref:     partial.NextRef,
					Factors: factors,
					Tuple:   newTuple,
				}

				err = detectLoop(partial, newPoint, traceable)
				if err != nil {
					results <- traceResult{error: fmt.Errorf("tracer detected a loop: %v", err)}
					return
				}

				builder := pathbuilder.Resume(partial.Path)

				if traceable.Segments() {
					builder.AddSegment()
				}

				builder.Add(newPoint)
				updatedPath := builder.Path()

				if traceable.Destination(job.destinationIPs, t.provider) {
					// Path is complete!
					results <- traceResult{path: &updatedPath}
					return
				}

				var wg sync.WaitGroup

				nextRefs, err := traceable.Next(newTuple, t.provider)
				if err != nil {
					results <- traceResult{error: fmt.Errorf("tracer was unable to get next refs: %v", err)}
					return
				}
				if len(nextRefs) < 1 {
					err := fmt.Errorf("no next points found when processing job:\n%+v", job)
					results <- traceResult{error: err}
					return
				}

				numNextRefs := len(nextRefs)
				resultChannels := make([]<-chan traceResult, numNextRefs)
				wg.Add(numNextRefs)
				for _, ref := range nextRefs {
					partial := reach.PartialPath{
						Path:    updatedPath,
						NextRef: ref,
					}
					j := tracerJob{
						partial:        partial,
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

func detectLoop(path reach.PartialPath, newPoint reach.Point, traceable reach.Traceable) error {
	if traceable.Visitable(path.Path.Contains(newPoint)) == false {
		return fmt.Errorf("cannot visit point again: %v", newPoint)
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
