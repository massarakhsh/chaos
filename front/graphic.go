package front

import (
	"github.com/massarakhsh/chaos/data"
)

type ItGraphic struct {
	ItPlot
}

func BuildGraphic(front *ItFront) *ItGraphic {
	it := &ItGraphic{}
	it.Front = front
	it.Loader = it
	it.IsZeroCenter = true
	it.Width, it.Height = 1024, 512
	return it
}

func (it *ItGraphic) Probe() bool {
	if dt := data.GetData(it.Sign, 4096); dt != nil {
		it.Load(dt)
		return true
	} else {
		return false
	}
}
