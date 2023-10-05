package pipeline

import (
	"context"
	"fmt"
)

// Builder contains context required by the pipeline package to construct a pipeline.
// Pipelines can be created using the New or NewFromChannel functions, and composed using the Connect function.
// A pipeline does not begin execution until Build is called on the stage.
type Builder[IN, OUT any] struct {
	prev       *Builder[any, IN]
	filters    FilterChain[IN]
	processor  Processor[IN, OUT]
	producer   Producer[OUT]
	bufferSize int
	workers    int
	jobs       []Job
}

// SetWorkerCount sets the number of goroutines which will run this stage of the pipeline.
func (b *Builder[IN, OUT]) SetWorkerCount(count int) {
	b.workers = count
}

// SetBufferSize causes this stage's output channel to be buffered with the specified capacity.
func (b *Builder[IN, OUT]) SetBufferSize(size int) {
	b.bufferSize = size
}

// FilterWith accepts one or more Filter functions and applies them to this stage of the pipeline BEFORE sending data
// to the stage's Processor.
func (b *Builder[IN, OUT]) FilterWith(f Filter[IN], fs ...Filter[IN]) *Builder[IN, OUT] {
	b.filters = append(b.filters, f)
	if len(fs) > 0 {
		b.filters = append(b.filters, fs...)
	}
	return b
}

// init recursively constructs the pipeline represented by a Builder. It returns a Pipeline struct, the output channel,
// and an error if anything failed during construction
func (b *Builder[IN, OUT]) init() (Pipeline, <-chan OUT, error) {
	if b.prev == nil {
		pipe := Pipeline{}
		if b.producer == nil {
			return pipe, nil, fmt.Errorf("init: first stage does not have a producer")
		}
		pipe.ctx, pipe.cancelFunc = context.WithCancel(context.Background())
		out, job := b.producer(pipe.ctx, b.bufferSize)
		pipe.jobs = append(pipe.jobs, job)
		return pipe, out, nil
	}

	pipe, in, err := b.prev.init()
	if err != nil {
		return pipe, nil, err
	}

	if b.processor == nil {
		return pipe, nil, fmt.Errorf("init: non-initial stage of pipeline does not have a processor")
	}

	out := make(chan OUT, b.bufferSize)
	for w := 0; w < b.workers; w++ {
		pipe.jobs = append(pipe.jobs, func() error {
			for {
				select {
				case <-pipe.ctx.Done():
					return ErrorCanceled
				case i, ok := <-in:
					if !ok {
						return nil
					}
					if b.filters.accept(i) {
						o, err := b.processor(i)
						if err != nil {
							return err
						}
						out <- o
					}
				}
			}
		})
	}

	return pipe, out, err
}

type FilterChain[T any] []Filter[T]

func (c FilterChain[T]) accept(item T) bool {
	for _, f := range c {
		if f(item) {
			return false
		}
	}
	return true
}

type Pipeline struct {
	ctx        context.Context
	cancelFunc func()
	waitFunc   func() error
	err        error
	stopped    bool
	jobs       []Job
}

func (p Pipeline) IsRunning() bool {
	return !p.stopped
}

func (p Pipeline) Cancel() {
	p.cancelFunc()
}

func (p Pipeline) Wait() error {
	if p.err == nil {
		p.err = p.waitFunc()
		p.stopped = true
	}
	return p.err
}
