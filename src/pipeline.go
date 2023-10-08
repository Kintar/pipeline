package pipeline

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"sync"
)

type JobSource interface {
	Previous() JobSource
	GetJob() Job
	GetWorkerCount() int
}

type ChannelSource[T any] interface {
	GetChannel() <-chan T
}

type stageBuilder[OUT any] struct {
	prev          JobSource
	workers       int
	buffer        int
	jobSource     func() Job
	channelSource func() chan OUT
}

func (s *stageBuilder[OUT]) GetWorkerCount() int {
	return s.workers
}

func (s *stageBuilder[OUT]) GetChannel() <-chan OUT {
	return s.channelSource()
}

func (s *stageBuilder[OUT]) SetWorkers(count int) {
	s.workers = count
}

func (s *stageBuilder[OUT]) SetBuffer(size int) {
	s.buffer = size
}

func (s *stageBuilder[OUT]) GetJob() Job {
	return s.jobSource()
}

func (s *stageBuilder[OUT]) Previous() JobSource {
	return s.prev
}

type Task struct {
	cancelFunc context.CancelCauseFunc
	waitFunc   Job
}

var ErrorCanceled = errors.New("pipeline canceled")

func (t *Task) Cancel() {
	t.cancelFunc(ErrorCanceled)
}

func (t *Task) Wait() error {
	return t.waitFunc()
}

func (s *stageBuilder[OUT]) Start() *Task {
	var prev JobSource
	ctx, cancelFunc := context.WithCancelCause(context.Background())
	g, _ := errgroup.WithContext(ctx)
	for prev = s; prev != nil; prev = prev.Previous() {
		g.Go(prev.GetJob())
	}
	return &Task{
		waitFunc:   g.Wait,
		cancelFunc: cancelFunc,
	}
}

type chn[T any] <-chan T

func (r chn[T]) GetChannel() <-chan T {
	return r
}

// New begins the construction of a pipeline by wrapping an existing channel in a ChannelSource that can be
// passed to Connect
func New[T any](in <-chan T) ChannelSource[T] {
	return chn[T](in)
}

// Connect performs the setup required to feed a Processor from a ChannelSource and use it as a Stage in a pipeline.
func Connect[IN, OUT any](input ChannelSource[IN], processor Processor[IN, OUT], filter PredicateFilter[OUT]) Stage[OUT] {
	s := &stageBuilder[OUT]{}

	s.channelSource = sync.OnceValue(func() chan OUT { return make(chan OUT, s.buffer) })
	s.jobSource = sync.OnceValue(func() Job {
		return makeProcessorPump[IN, OUT](input.GetChannel(), s.channelSource(), processor, filter)
	})

	return s
}

type Stage[OUT any] interface {
	Previous() JobSource
	GetChannel() <-chan OUT
	SetWorkers(count int)
	SetBuffer(size int)
	GetJob() Job
}
