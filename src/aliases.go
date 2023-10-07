package pipeline

// Processor is a function that accepts data of one type, performs work, and returns data of a different type.
type Processor[IN, OUT any] func(IN) (OUT, error)

// Predicate is used to perform a logic test on an item.
type Predicate[T any] func(T) bool

type PredicateFilter[T any] []Predicate[T]

// Accept returns true if all Predicates return true, otherwise false
func (f PredicateFilter[T]) Accept(item T) bool {
	for _, p := range f {
		if !p(item) {
			return false
		}
	}
	return true
}

type Job func() error
