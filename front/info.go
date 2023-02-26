package front

import (
	"github.com/massarakhsh/chaos/pkg/zone"
)

type ItInfo struct {
	zone.ItZone
	viewSign int
}

func buildInfo() *ItInfo {
	it := &ItInfo{}
	it.BindVerticalBox(it)
	it.buildInfo()
	return it
}

func (it *ItInfo) buildInfo() {
	gra := BuildGraphic()
	it.Append(gra, true)
	if down := zone.BuildHorizontalBox(nil); down != nil {
		it.Append(down, true)
		if child := BuildSpectr(); child != nil {
			down.Append(child, true)
		}
	}
}
