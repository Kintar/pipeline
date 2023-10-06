package pipeline

import (
	"github.com/Kintar/pipeline/funcs"
	"testing"
)

func Test_Connections(t *testing.T) {
	var src Emitter[int]
	var proc funcs.Processor[int, float32]
	var proc2 funcs.Processor[float32, string]
	var f PredicateFilter[int]

	t1 := Begin(src, proc)

	Connect(
		t1,
		proc2)
}
