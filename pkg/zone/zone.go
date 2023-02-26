package zone

import (
	"github.com/andlabs/ui"
)

type ItZone struct {
	self     ItfZone
	refresh  ItfRefresh
	box      *ui.Box
	ctrl     ui.Control
	listZone []ItfZone
}

type ItfZone interface {
	BindHorizontalBox(self ItfZone)
	BindVerticalBox(self ItfZone)
	BindControl(self ItfZone, ctrl ui.Control)
	BindRefresh(refresh ItfRefresh)
	GetControl() ui.Control
	Append(child ItfZone, stretchy bool)
	GetVolume() int
	PushControl(ctrl ui.Control)
	PopControls(last int)
	Step()
}

type ItfRefresh interface {
	Refresh()
}

func BuildHorizontalBox(refresh ItfRefresh) ItfZone {
	it := &ItZone{}
	it.BindHorizontalBox(it)
	it.BindRefresh(refresh)
	return it
}

func BuildVerticalBox(refresh ItfRefresh) ItfZone {
	it := &ItZone{}
	it.BindVerticalBox(it)
	it.BindRefresh(refresh)
	return it
}

func BuildFromControl(ctrl ui.Control) ItfZone {
	it := &ItZone{}
	it.BindControl(it, ctrl)
	return it
}

func (it *ItZone) BindHorizontalBox(self ItfZone) {
	it.self = self
	it.box = ui.NewHorizontalBox()
	it.box.SetPadded(true)
}

func (it *ItZone) BindVerticalBox(self ItfZone) {
	it.self = self
	it.box = ui.NewVerticalBox()
	it.box.SetPadded(true)
}

func (it *ItZone) BindControl(self ItfZone, ctrl ui.Control) {
	it.self = self
	it.ctrl = ctrl
}

func (it *ItZone) BindRefresh(refresh ItfRefresh) {
	it.refresh = refresh
}

func (it *ItZone) GetVolume() int {
	return len(it.listZone)
}

func (it *ItZone) GetControl() ui.Control {
	if it.box != nil {
		return it.box
	} else {
		return it.ctrl
	}
}

func (it *ItZone) Append(child ItfZone, stretchy bool) {
	if it.box != nil {
		if ctrl := child.GetControl(); ctrl != nil {
			it.box.Append(ctrl, stretchy)
			it.listZone = append(it.listZone, child)
		}
	}
}

func (it *ItZone) PushControl(ctrl ui.Control) {
	if child := BuildFromControl(ctrl); child != nil {
		it.Append(child, false)
	}
}

func (it *ItZone) PopControls(last int) {
	if it.box != nil {
		for pos := len(it.listZone) - 1; pos > last; pos-- {
			it.box.Delete(pos)
			it.listZone = it.listZone[:pos]
		}
	}
}

func (it *ItZone) Step() {
	if it.refresh != nil {
		it.refresh.Refresh()
	}
	for _, zone := range it.listZone {
		zone.Step()
	}
}
