package front

import (
	"github.com/andlabs/ui"
)

type ItImage struct {
	ItPanel
}

func (it *ItImage) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	it.resize(p)
	//it.clear(p)
	//it.calc(p)
	//it.drawAxes(p)
	//it.drawGraph(p)
	//it.drawPens(p)
	//img := ui.NewImage(it.Width, it.Height)
}

func (it *ItImage) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	/*xs, ys := pointLocations(it.Width, it.Height)

	data.Ser.Current = -1
	for i := 0; i < len(xs); i++ {
		if inPoint(me.X, me.Y, xs[i], ys[i]) {
			data.Ser.Current = i
			break
		}
	}*/

	a.QueueRedrawAll()
}

func (it *ItImage) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (it *ItImage) DragBroken(a *ui.Area) {
	// do nothing
}

func (it *ItImage) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}
