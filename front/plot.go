package front

import (
	"fmt"
	"time"

	"github.com/andlabs/ui"
)

const bd = 10

//const radius = 5

type ItPlot struct {
	ItPanel
	SP         *ui.DrawStrokeParams
	area       *ui.Area
	lastUpdate time.Time
}

type ItPoint struct {
	XVal, YVal float64
	XLoc, YLoc float64
}

func (it *ItPlot) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	//it.Loader.Probe()
	if it.Name == "signal" {
		it.Name += ""
	}
	it.resize(p)
	it.clear(p)
	it.calc(p)
	it.drawAxes(p)
	it.drawGraph(p)
	it.drawPens(p)
}

func (it *ItPlot) clear(p *ui.AreaDrawParams) {
	it.SP = &ui.DrawStrokeParams{
		Cap:        ui.DrawLineCapFlat,
		Join:       ui.DrawLineJoinMiter,
		Thickness:  2,
		MiterLimit: ui.DrawDefaultMiterLimit,
	}
	brush := mkSolidBrush(0xffffff, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(0, 0, it.Width, it.Height)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()
}

func (it *ItPlot) calc(p *ui.AreaDrawParams) {
	it.X.Calibrate(bd, it.Width-2*bd, true)
	it.Y.Calibrate(bd, it.Height-2*bd, false)
	it.X.LocZero = it.X.LocDep
	it.Y.LocZero = it.Y.LocDep
	if length := it.Count; length >= 2 {
		for i := 0; i < length; i++ {
			point := &it.List[i]
			point.XLoc = it.X.ToLoc(point.XVal)
			point.YLoc = it.Y.ToLoc(point.YVal)
		}
	}
}

func (it *ItPlot) drawAxes(p *ui.AreaDrawParams) {
	if path := ui.DrawNewPath(ui.DrawFillModeWinding); path != nil {
		for x := it.X.First; x < it.X.Max; x += it.X.Step {
			if xt := it.X.ToLoc(x); xt >= 0 && xt < it.Width {
				path.NewFigure(it.X.LocDep+xt, it.X.LocDep+it.Y.LocSize)
				path.LineTo(it.X.LocDep+xt, it.X.LocDep)
				it.drawText(p, fmt.Sprintf(it.X.Format, x), it.X.LocDep+xt-16, it.Y.LocDep+it.Y.LocSize-10, 10, 0.99, 0, 0, 0.99)
			}
		}
		path.End()
		brush := mkSolidBrush(0x808080, 0.5)
		p.Context.Stroke(path, brush, it.SP)
		path.Free()
	}
	if path := ui.DrawNewPath(ui.DrawFillModeWinding); path != nil {
		path.NewFigure(it.X.LocZero, it.Y.LocDep+it.Y.LocSize)
		path.LineTo(it.X.LocZero, it.Y.LocDep)
		path.NewFigure(it.X.LocDep, bd+it.Y.LocSize-it.Y.LocZero)
		path.LineTo(it.X.LocDep+it.X.LocSize, bd+it.Y.LocSize-it.Y.LocZero)
		path.End()
		brush := mkSolidBrush(0x000000, 1.0)
		p.Context.Stroke(path, brush, it.SP)
		path.Free()
	}
}

func (it *ItPlot) drawText(p *ui.AreaDrawParams, text string, x, y float64, size float64, cr, cg, cb, ca float64) {
	attrstr := ui.NewAttributedString(text)
	sz := len(text)
	attrstr.SetAttribute(ui.TextColor{R: cr, G: cg, B: cb, A: ca}, 0, sz)
	df := ui.FontDescriptor{
		Family:  "Courier New",
		Size:    ui.TextSize(size),
		Weight:  ui.TextWeightNormal,
		Italic:  ui.TextItalicNormal,
		Stretch: ui.TextStretchCondensed,
	}
	lp := ui.DrawTextLayoutParams{
		String:      attrstr,
		DefaultFont: &df,
		Width:       it.Width - x,
		Align:       ui.DrawTextAlignLeft,
	}
	tl := ui.DrawNewTextLayout(&lp)
	p.Context.Text(tl, x, y)
	tl.Free()
}

func (it *ItPlot) drawGraph(p *ui.AreaDrawParams) {
	brush := mkSolidBrush(0x5599ff, 0.5)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	if len(it.List) > 0 {
		path.NewFigure(it.List[0].XLoc, it.List[0].YLoc)
		for i := 1; i < it.Count; i++ {
			path.LineTo(it.List[i].XLoc, it.List[i].YLoc)
		}
	}
	path.End()
	p.Context.Stroke(path, brush, it.SP)
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

func (it *ItPlot) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
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

func (it *ItPlot) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (it *ItPlot) DragBroken(a *ui.Area) {
	// do nothing
}

func (it *ItPlot) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}

func mkSolidBrush(color uint32, alpha float64) *ui.DrawBrush {
	brush := new(ui.DrawBrush)
	brush.Type = ui.DrawBrushTypeSolid
	component := uint8((color >> 16) & 0xFF)
	brush.R = float64(component) / 255
	component = uint8((color >> 8) & 0xFF)
	brush.G = float64(component) / 255
	component = uint8(color & 0xFF)
	brush.B = float64(component) / 255
	brush.A = alpha
	return brush
}
