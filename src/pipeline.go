package pipeline

func new[IN, OUT any](input <-chan IN, proc Processor[IN, OUT]) stage[OUT] {
	
}
