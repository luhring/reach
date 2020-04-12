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

func (t *Tracer) Trace(source, destination reach.Subject) []reach.Path {
	sourceRef := reach.InfrastructureReference{
		R: reach.ResourceReference{
			Domain: source.Domain,
			Kind:   source.Kind,
			ID:     source.ID,
		},
		Implicit: false,
	}

	firstPartialPath := reach.PartialPath{
		Path:    reach.NewPath(),
		NextRef: sourceRef,
	}
	firstJob := tracerJob{
		source:      source,
		destination: destination,
		partialPath: firstPartialPath,
	}

	done := make(chan interface{})
	defer close(done)

	results := t.tracePoint(done, firstJob)
	var paths []reach.Path

	for result := range results {
		if result.error != nil {
			_, _ = fmt.Fprintln(os.Stderr, result.error) // TODO: Log more intelligently!
		}
		paths = append(paths, *result.complete)
	}

	return paths
}

type result struct {
	complete *reach.Path
	error    error
}

func (t *Tracer) tracePoint(done <-chan interface{}, job tracerJob) <-chan result {
	results := make(chan result)

	go func() {
		defer close(results)
		for {
			select {
			case <-done:
				return
			default:
				current := job.partialPath
				ref := current.NextRef

				// Turn ref into full TraceableInfrastructure
				r, err := t.provider.Get(ref)
				if err != nil {
					results <- result{error: err}
				}
				infra := r.Properties.(reach.TraceableInfrastructure) // TODO: Make this more type-safe

				factors := infra.Factors()
				prevTuple := current.Path.LastPoint().Tuple
				newTuple := infra.UpdatedTuple(&prevTuple)

				newPoint := reach.Point{
					Ref:     ref,
					Factors: factors,
					Tuple:   newTuple,
				}

				builder := reach.ResumePathBuilding(current.Path)

				if infra.Segments() {
					builder.AddSegment()
				}

				builder.Add(newPoint)
				updatedPath := builder.Path()

				var destIPs []net.IP // TODO: Get from destination subject
				if infra.IsDestinationForIP(destIPs) {
					// Path is complete!
					results <- result{complete: &updatedPath}
					return
				}

				var wg sync.WaitGroup

				nextRefs := infra.Next(newTuple)
				if len(nextRefs) < 1 {
					err := fmt.Errorf("no next points found when processing job:\n%+v", job)
					results <- result{
						error: err,
					}
				}

				numNextRefs := len(nextRefs)
				resultChannels := make([]<-chan result, numNextRefs)
				wg.Add(numNextRefs)
				for _, ref := range nextRefs { // TODO: Make sure these can get fired off! Not block! Parent goroutine can stay open, though.
					partial := reach.PartialPath{
						Path:    updatedPath,
						NextRef: ref,
					}
					j := tracerJob{
						source:      job.source,
						destination: job.destination,
						partialPath: partial,
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

				return
			}
		}
	}()

	return results
}

func fanIn(
	done <-chan interface{},
	channels []<-chan result,
) <-chan result {
	var wg sync.WaitGroup
	multiplexedStream := make(chan result)

	multiplex := func(c <-chan result) {
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
