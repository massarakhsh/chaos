package front

import (
	"github.com/andlabs/ui"
)

const bd = 10
const radius = 5

type ItSource interface {
	GetData() *ItData
	LoadData()
}

type ItPlot struct {
	Source ItSource
}

func (it ItPlot) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	it.Source.LoadData()
	it.resize(p)
	it.clear(p)
	it.calc(p)
	it.drawAxes(p)
	it.drawGraph(p)
	it.drawPens(p)
}

func (it ItPlot) resize(p *ui.AreaDrawParams) {
	data := it.Source.GetData()
	data.Width, data.Height = p.AreaWidth-2*bd, p.AreaHeight-2*bd
	m := ui.DrawNewMatrix()
	m.Translate(bd, bd)
	p.Context.Transform(m)
}

func (it ItPlot) clear(p *ui.AreaDrawParams) {
	data := it.Source.GetData()
	data.SP = &ui.DrawStrokeParams{
		Cap:        ui.DrawLineCapFlat,
		Join:       ui.DrawLineJoinMiter,
		Thickness:  2,
		MiterLimit: ui.DrawDefaultMiterLimit,
	}
	brush := mkSolidBrush(colorWhite, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(-bd, -bd, data.Width+2*bd, data.Height+2*bd)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()
}

func (it ItPlot) calc(p *ui.AreaDrawParams) {
	data := it.Source.GetData()
	data.Calc()
}

func (it ItPlot) drawAxes(p *ui.AreaDrawParams) {
	data := it.Source.GetData()
	brush := mkSolidBrush(colorBlack, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigure(data.XZero, 0)
	path.LineTo(data.XZero, data.Height)
	path.NewFigure(0, data.YZero)
	path.LineTo(data.Width, data.YZero)
	path.End()
	p.Context.Stroke(path, brush, data.SP)
	path.Free()
}

func (it *ItPlot) drawGraph(p *ui.AreaDrawParams) {
	data := it.Source.GetData()
	brush := mkSolidBrush(colorDodgerBlue, 0.5)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigure(data.List[0].XLoc, data.List[0].YLoc)
	for i := 1; i < data.Count; i++ {
		path.LineTo(data.List[i].XLoc, data.List[i].YLoc)
	}
	path.End()
	p.Context.Stroke(path, brush, data.SP)
	path.Free()
}

func (it *ItPlot) drawPens(p *ui.AreaDrawParams) {
	/*if data.Ser.Current >= 0 {
		xs, ys := pointLocations(it.Width, it.Height)
		path := ui.DrawNewPath(ui.DrawFillModeWinding)
		path.NewFigureWithArc(
			xs[data.Ser.Current], ys[data.Ser.Current],
			radius,
			0, 6.23, // TODO pi
			false)
		path.End()
		// use the same brush as for the histogram lines
		brush := mkSolidBrush(colorDodgerBlue, 1.0)
		p.Context.Fill(path, brush)
		path.Free()
	}*/
}

func (it ItPlot) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
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

func (it ItPlot) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (it ItPlot) DragBroken(a *ui.Area) {
	// do nothing
}

func (it ItPlot) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}
