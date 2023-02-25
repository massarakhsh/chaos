package front

import (
	"github.com/massarakhsh/chaos/data"
)

type ItSignal struct {
	ItPlot
}

func BuildSignal(front *ItFront) *ItSignal {
	it := &ItSignal{}
	it.Front = front
	it.Name = "signal"
	it.Loader = it
	it.IsZeroCenter = true
	it.Width, it.Height = 4096, 1024
	return it
}

func (it *ItSignal) Probe() bool {
	if dt := data.GetData(it.Sign, 4096); dt != nil {
		it.Load(dt)
		return true
	} else {
		return false
	}
}
