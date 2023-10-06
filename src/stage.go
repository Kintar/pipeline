package pipeline

import (
	"context"
	"github.com/Kintar/pipeline/funcs"
)

type TaskDef[IN, OUT any] struct {
	filter       PredicateFilter[IN]
	processor    funcs.Processor[IN, OUT]
	outputBuffer int
	workerCount  int
}

func (s *Stage[IN, OUT]) AsEmitter() Emitter[OUT] {
	return s
}

func (t *TaskDef[IN, OUT]) Workers(count int) *TaskDef[IN, OUT] {
	t.workerCount = count
	return t
}

func (t *TaskDef[IN, OUT]) Buffer(size int) *TaskDef[IN, OUT] {
	t.outputBuffer = size
	return t
}

func (t *TaskDef[IN, OUT]) Filters(fs ...funcs.Predicate[IN]) *TaskDef[IN, OUT] {
	t.filter = fs
	return t
}

func (t *TaskDef[IN, OUT]) Filter(f PredicateFilter[IN]) *TaskDef[IN, OUT] {
	t.filter = f
	return t
}

type Stage[IN, OUT any] struct {
	*TaskDef[IN, OUT]
	ctx    context.Context
	prev   Emitter[IN]
	output chan OUT
	pump   funcs.Job
}

func (s *Stage[IN, OUT]) Processor() funcs.Processor[IN, OUT] {
	return s.processor
}

func (s *Stage[IN, OUT]) getEmitter() <-chan OUT {
	if s.output == nil {
		s.output = make(chan OUT, s.outputBuffer)

		if len(s.filter) == 0 {
			s.pump = processOnly(s.prev.getContext(), s.prev.getEmitter(), s.output, s.processor)
		} else {
			s.pump = filterAndProcess(s.prev.getContext(), s.prev.getEmitter(), s.output, s.filter, s.processor)
		}
	}

	return s.output
}

func (s *Stage[IN, OUT]) getContext() context.Context {
	return s.prev.getContext()
}

func (s *Stage[IN, OUT]) receiveFrom(input Emitter[IN]) {
	s.prev = input
}

func filterAndProcess[IN, OUT any](ctx context.Context, input <-chan IN, output chan<- OUT, pred PredicateFilter[IN], proc funcs.Processor[IN, OUT]) funcs.Job {
	return func() error {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return ErrorCanceled
			case item, ok := <-input:
				if !ok {
					return nil
				}
				if !pred.Accept(item) {
					continue
				}
				result, err := proc(item)
				if err != nil {
					return err
				}
				output <- result
			}
		}
	}
}

func processOnly[IN, OUT any](ctx context.Context, input <-chan IN, output chan<- OUT, proc funcs.Processor[IN, OUT]) funcs.Job {
	return func() error {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return ErrorCanceled
			case item, ok := <-input:
				if !ok {
					return nil
				}
				result, err := proc(item)
				if err != nil {
					return err
				}
				output <- result
			}
		}
	}
}
