package pipeline

type stage[OUT any] struct {
	jobFunc  func() Job
	supplier channelSupplier[OUT]
}

func (s stage[OUT]) supplyChannel() chan OUT {
	return s.supplier.supplyChannel()
}

func (s stage[OUT]) close() {
	s.supplier.close()
}

func (s stage[OUT]) getJob() Job {
	return s.jobFunc()
}
