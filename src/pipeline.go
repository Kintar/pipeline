package pipeline

import "github.com/Kintar/pipeline/funcs"

type receiverBuilder[IN any] interface {
	input() chan<- IN
}

type emitterBuilder[OUT any] interface {
	output() <-chan OUT
}

// Pipe represents the connection between an Emitter and a Receiver
type Pipe[IN, OUT any] struct {
	in  emitterBuilder[IN]
	out receiverBuilder[OUT]
}

type Stage[IN, OUT any] struct {
	source    funcs.Producer[IN]
	filter    PredicateFilter[IN]
	processor funcs.Processor[IN, OUT]
}

func buildPipe[T, U, IN, OUT any](s1 Stage[IN, OUT])
