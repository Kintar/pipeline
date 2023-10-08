package pipeline

import (
	"sync"
)

func makeProcessorPump[IN, OUT any](
	input <-chan IN,
	output chan<- OUT,
	proc Processor[IN, OUT],
	filter PredicateFilter[OUT],
) Job {
	return func() error {
		defer sync.OnceFunc(func() { close(output) })
		out := output
		for i := range input {
			o, err := proc(i)
			if err != nil {
				return err
			}
			if filter.Accept(o) {
				out <- o
			}
		}
		return nil
	}
}
