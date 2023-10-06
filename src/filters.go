package pipeline

import "github.com/Kintar/pipeline/funcs"

type PredicateFilter[T any] []funcs.Predicate[T]

// Accept returns true if all Predicates return true, otherwise false
func (f PredicateFilter[T]) Accept(item T) bool {
	for _, p := range f {
		if !p(item) {
			return false
		}
	}
	return true
}
