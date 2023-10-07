package pipeline

type channelSupplier[T any] interface {
	supplyChannel() chan T
	close()
}

type builder[IN, OUT any] struct {
	input        channelSupplier[IN]
	supplierFunc func() channelSupplier[OUT]
	processor    Processor[IN, OUT]
	filter       PredicateFilter[OUT]
	workerCount  int
	bufferSize   int
}

func (b *builder[IN, OUT]) supplyChannel() chan OUT {
	return b.supplierFunc().supplyChannel()
}

func (b *builder[IN, OUT]) close() {
	b.supplierFunc().close()
}

func (b *builder[IN, OUT]) Filter(f PredicateFilter[OUT]) *builder[IN, OUT] {
	b.filter = f
	return b
}

func (b *builder[IN, OUT]) SetBufferSize(size int) *builder[IN, OUT] {
	b.bufferSize = size
	return b
}

func (b *builder[IN, OUT]) SetWorkerCount(count int) *builder[IN, OUT] {
	b.workerCount = count
	return b
}

func (b *builder[IN, OUT]) Build() stage[OUT] {
	s := stage[OUT]{
		supplier: b,
	}

	job := makeProcessorPump(b.input, s.supplier, b.processor, b.filter)
}
