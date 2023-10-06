package tmp

import "context"

type Emitter[T any] interface {
	getEmitter() <-chan T
	getContext() context.Context
}

type Processor[IN, OUT any] func(IN) (OUT, error)

type TaskDef[IN, OUT any] struct {
	// TaskDef is composed of its processor and various other things
	Processor[IN, OUT]
	// redacted for space
}

// implementes Emitter[OUT] for the TaskDef
func (t *TaskDef[IN, OUT]) getEmitter() <-chan OUT {
	// redacted for space
	panic("not implemented")
}

func (t *TaskDef[IN, OUT]) getContext() context.Context {
	/// redacted for space
	panic("not implemented")
}

func Begin[IN, OUT any](in Emitter[IN], proc Processor[IN, OUT]) *TaskDef[IN, OUT] {
	// redacted for space
	panic("not implemented")
}

func Connect[T, IN, OUT any](in *TaskDef[T, IN], proc Processor[IN, OUT]) *TaskDef[IN, OUT] {
	panic("not implemented")
}

func example() {
	var source Emitter[int]
	var p1 Processor[int, float32]

	// initialize variables, obviously

	stage1 := Begin(source, p1)

	var p2 Processor[float32, string]
	stage2 := Connect(stage1, p2)
}
