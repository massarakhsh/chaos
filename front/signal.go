package front

import (
	"github.com/massarakhsh/chaos/data"
)

type ItSignal struct {
	ItImage
}

func BuildSignal() *ItSignal {
	graph := &ItSignal{}
	graph.Loader = graph
	graph.IsZeroCenter = true
	graph.Width, graph.Height = 1024, 512
	return graph
}

func (it *ItSignal) Probe() bool {
	if dt := data.GetData(it.Sign, 4096); dt != nil {
		it.Load(dt)
		return true
	} else {
		return false
	}
}
