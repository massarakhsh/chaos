package front

import (
	"github.com/massarakhsh/chaos/pkg/zone"
)

type ItInfo struct {
	zone.ItZone
}

func buildInfo() *ItInfo {
	it := &ItInfo{}
	it.BindVerticalBox(it)
	it.buildInfo()
	return it
}

func (it *ItInfo) buildInfo() {
	gra := BuildSignal()
	it.Append(gra, true)
	inter := BuildInterval(&gra.Panel)
	it.Append(inter, true)
	down := zone.BuildHorizontalBox(nil)
	it.Append(down, true)
	child := BuildSpectr(inter)
	down.Append(child, true)
	crop := BuildInterval(&child.Panel)
	down.Append(crop, true)
}
