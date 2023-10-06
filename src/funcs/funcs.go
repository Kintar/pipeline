package funcs

// Job is a long-running function which returns nil to indicate successful completion, or an error to indicate failure.
type Job func() error

// Producer is a function that returns a series of items. If no more values are available, the Producer should return
// the last item in the series and false. A non-nil error indicates that there was an issue producing the item.
//
// Special care must be taken if the producer can fail to produce even a single item. In these cases, the calling func
// must know how to interpret the difference between the final call of a Producer which is returning a copy of the final
// element and false, and the initial call of a Producer which is returning a zero value and false.
type Producer[T any] func() (T, bool, error)

// Processor is a function that accepts data of one type, performs work, and returns data of a different type.
type Processor[IN, OUT any] func(IN) (OUT, error)

// Consumer is a function that accepts an item and performs work.
type Consumer[IN any] func(IN) error

// Predicate is used to perform a logic test on an item.
type Predicate[T any] func(T) bool
