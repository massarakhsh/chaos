package front

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

type ItOrg struct {
}

var MainOrg ItOrg

func BuildOrg() *ItOrg {
	return &MainOrg
}

func (it *ItOrg) Probe() bool {
	return false
}

func (it *ItOrg) Draw(a *ui.Area, p *ui.AreaDrawParams) {
}

func (it *ItOrg) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	a.QueueRedrawAll()
}

func (it *ItOrg) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (it *ItOrg) DragBroken(a *ui.Area) {
	// do nothing
}

func (it *ItOrg) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}
