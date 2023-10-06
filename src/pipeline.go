package pipeline

import (
	"context"
	"errors"
	"github.com/Kintar/pipeline/funcs"
)

/*
source := pipeline.FromProducer(prod).Buffer(200)
task1 := pipeline.Connect(

	source,
	pipeline.Task(myProcessorFunc, filters).Workers(5).Buffer(0),

)

task2 := pipeline.Connect(

	task1,
	pipeline.Task(otherProcessor).Workers(5),

)

job := pipeline.Consume(task2, consumerThing)

output, job := pipeline.Emit(task2).Workers(2).Buffer(0)
*/

var ErrorCanceled = errors.New("pipeline aborted: context cancelled")

type Emitter[T any] interface {
	getEmitter() <-chan T
	getContext() context.Context
}

func FromProducer[T any](p funcs.Producer[T]) Emitter[T] {
	panic("implement me")
}

func Task[IN, OUT any](processor funcs.Processor[IN, OUT]) *TaskDef[IN, OUT] {
	t := &TaskDef[IN, OUT]{
		processor:   processor,
		workerCount: 1,
	}

	return t
}

func Begin[IN, OUT any](in Emitter[IN], p funcs.Processor[IN, OUT]) *Stage[IN, OUT] {
	panic("not implemented")
}

func Connect[T, IN, OUT any](in *Stage[T, IN], p funcs.Processor[IN, OUT]) *Stage[IN, OUT] {
	s := &Stage[IN, OUT]{
		prev: in,
		TaskDef: &TaskDef[IN, OUT]{
			processor: p,
		},
	}

	return s
}
