package pipeline

func makeProcessorPump[IN, OUT any](
	input channelSupplier[IN],
	output channelSupplier[OUT],
	proc Processor[IN, OUT],
	filter PredicateFilter[OUT],
) Job {
	return func() error {
		out := output.supplyChannel()
		for i := range input.supplyChannel() {
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
