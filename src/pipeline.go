package pipeline

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
)

// ErrorCanceled is returned by a pipeline's Wait method if the pipeline ended due to cancellation of the controlling
// Context
var ErrorCanceled = errors.New("pipeline aborted: context canceled")

// jobNoOp is a known Job that does...nothing. It simply returns nil.
var jobNoOp = func() error { return nil }

// Job is a long-running function which returns nil to indicate successful completion, or an error to indicate abnormal
// completion. Within the context of a Pipeline, a Job which returns an error will cause the entire pipeline to be
// canceled.
type Job func() error

// Source is a function that returns a series of items. If more items are available, the boolean value must be true.
// If there was an unrecoverable error that causes the producer to fail before reaching the end of its sequence, the
// error value must be non-nil.
type Source[T any] func() (T, bool, error)

// Producer is a function that generates a channel and a job to feed that channel.
type Producer[T any] func(ctx context.Context, bufferSize int) (<-chan T, Job)

// Processor is a function that accepts data of one type, performs work, and returns data of a different type.
// If the Processor returns a non-nil error value, its associated pipeline will be aborted with the error as the reason.
type Processor[IN, OUT any] func(IN) (OUT, error)

// Consumer is a function that accepts the final data type in a pipeline and performs some work. If a Consumer returns
// a non-nil error, its associated pipeline is aborted with the error as the reason
type Consumer[IN any] func(IN) error

// Filter is used to discard some elements of a pipeline from further processing. If a filter returns true, the element
// is dropped and not sent further down the pipeline.
type Filter[T any] func(T) bool

// channelProducer wraps an existing channel into a Producer of its type. When the Producer is executed, if the requested
// bufferSize is larger than the input channel, a new channel with the necessary size is allocated and the resulting Job
// will feed data into the new channel until the original is closed. Otherwise, the original channel is returned and the
// Job does nothing.
func channelProducer[T any](input <-chan T) Producer[T] {
	return func(ctx context.Context, bufferSize int) (<-chan T, Job) {
		// If the requested buffer is larger than our existing channel, make a new channel and a Job to populate it
		if bufferSize < cap(input) {
			buffered := make(chan T, bufferSize)
			return buffered, func() error {
				defer close(buffered)
				for {
					select {
					case <-ctx.Done():
						return ErrorCanceled
					case i, ok := <-input:
						if !ok {
							return nil
						}
						buffered <- i
					}
				}
			}
		}

		// Otherwise, we can just return the input and a no-op job
		return input, jobNoOp
	}
}

func sourceProducer[T any](s Source[T]) Producer[T] {
	return func(ctx context.Context, bufferSize int) (<-chan T, Job) {
		outChan := make(chan T, bufferSize)
		return outChan, func() error {
			defer close(outChan)
			for {
				item, more, err := s()
				if err != nil {
					return err
				}
				outChan <- item
				if !more {
					return nil
				}
			}
		}
	}
}

// New begins a new pipeline with the given Source as its source.
func New[T any](s Source[T]) *Builder[T, T] {
	return NewFromProducer(sourceProducer(s))
}

// NewFromProducer begins a new pipeline with the given Producer as its source.
func NewFromProducer[T any](p Producer[T]) *Builder[T, T] {
	return &Builder[T, T]{
		producer:   p,
		bufferSize: 0,
		workers:    1,
	}
}

// NewFromChannel begins a new pipeline with the given channel as its source. When the channel is closed, the pipeline
// is concluded.
func NewFromChannel[T any](c <-chan T) *Builder[T, T] {
	return &Builder[T, T]{
		producer:   channelProducer[T](c),
		bufferSize: 0,
		workers:    1,
	}
}

// Connect extends the given Builder and sends its output to the given processor and returns a new Builder
func Connect[IN, OUT any](b *Builder[any, IN], processor Processor[IN, OUT]) *Builder[IN, OUT] {
	return &Builder[IN, OUT]{
		prev:       b,
		processor:  processor,
		bufferSize: 0,
		workers:    1,
	}
}

// Emit builds the pipeline represented by Builder and returns the pipeline and a channel which will be closed when the
// pipeline completes or errors.
func Emit[IN, OUT any](b *Builder[IN, OUT]) (Pipeline, <-chan OUT, error) {
	pipe, out, err := b.init()
	if err != nil {
		return pipe, out, err
	}
	g, _ := errgroup.WithContext(pipe.ctx)
	for _, j := range pipe.jobs {
		g.Go(j)
	}
	pipe.waitFunc = g.Wait
	return pipe, out, err
}

// Consume builds the pipeline represented by Builder and sends the output to the specified Consumer.
func Consume[OUT any](b *Builder[any, OUT], consumer Consumer[OUT]) (Pipeline, error) {
	pipe, out, err := b.init()
	if err != nil {
		return pipe, err
	}
	g, _ := errgroup.WithContext(pipe.ctx)
	for _, j := range pipe.jobs {
		g.Go(j)
	}
	pipe.waitFunc = g.Wait
	g.Go(func() error {
		for i := range out {
			err := consumer(i)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return pipe, nil
}
