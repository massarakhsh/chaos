package front

import (
	"fmt"
	"math"

	"github.com/andlabs/ui"
)

const bd = 10

//const radius = 5

type ItPlot struct {
	ItPanel
	SP *ui.DrawStrokeParams
}

type InfPlot interface {
	ui.AreaHandler
	Probe() bool
}

type ItPoint struct {
	XVal, YVal float64
	XLoc, YLoc float64
}

func (it *ItPlot) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	//it.Loader.Probe()
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
	path.AddRectangle(-bd, -bd, it.Width+2*bd, it.Height+2*bd)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()
}

func (it *ItPlot) calc(p *ui.AreaDrawParams) {
	it.XZero, it.YZero = 0, it.Height
	if length := it.Count; length >= 2 {
		it.XZero = it.locFromX(0)
		if it.XZero < 0 {
			it.XZero = 0
		} else if it.XZero > it.Width {
			it.XZero = it.Width
		}
		it.YZero = it.locFromY(0)
		if it.YZero < 0 {
			it.YZero = 0
		} else if it.YZero > it.Height {
			it.YZero = it.Height
		}
		for i := 0; i < length; i++ {
			point := &it.List[i]
			point.XLoc = it.locFromX(point.XVal)
			point.YLoc = it.locFromY(point.YVal)
		}
		it.XFirst, it.XStep = it.findScale(it.XMin, it.XMax)
		it.YFirst, it.YStep = it.findScale(it.YMin, it.YMax)
	}
}

func (it *ItPlot) locFromX(x float64) float64 {
	return (0*(it.XMax-x) + it.Width*(x-it.XMin)) / (it.XMax - it.XMin)
}

func (it *ItPlot) locFromY(y float64) float64 {
	return (0*(y-it.YMin) + it.Height*(it.YMax-y)) / (it.YMax - it.YMin)
}

func (it *ItPlot) locToX(x float64) float64 {
	return (it.XMin*(it.Width-x) + it.XMax*(x-0)) / it.Width
}

func (it *ItPlot) locToY(y float64) float64 {
	return (it.YMin*(y-0) + it.YMax*(it.Height-y)) / it.Height
}

func (it *ItPlot) findScale(min, max float64) (float64, float64) {
	step := 1.0
	for step*10 > max-min {
		step /= 10
	}
	for step*100 < max-min {
		step *= 10
	}
	for step*20 < max-min {
		step *= 2
	}
	first := math.Floor(min/step) - step
	for first < min {
		first += step
	}
	return first, step
}

func (it *ItPlot) drawAxes(p *ui.AreaDrawParams) {
	if path := ui.DrawNewPath(ui.DrawFillModeWinding); path != nil {
		y := 0.0
		for x := it.XFirst; x < it.XMax; x += it.XStep {
			if xt := it.locFromX(x); xt >= 0 && xt < it.Width {
				path.NewFigure(xt, 0)
				path.LineTo(xt, it.Height)
				it.drawText(p, fmt.Sprintf("%f", x), xt, y, 10, 0.99, 0, 0, 0.99)
				y += 12
			}
		}
		path.End()
		brush := mkSolidBrush(0x808080, 0.5)
		p.Context.Stroke(path, brush, it.SP)
		path.Free()
	}
	if path := ui.DrawNewPath(ui.DrawFillModeWinding); path != nil {
		path.NewFigure(it.XZero, 0)
		path.LineTo(it.XZero, it.Height)
		path.NewFigure(0, it.YZero)
		path.LineTo(it.Width, it.YZero)
		path.End()
		brush := mkSolidBrush(0x000000, 1.0)
		p.Context.Stroke(path, brush, it.SP)
		path.Free()
	}
}

/*func (it *ItPlot) appendWithAttributes(what string, attrs ...ui.Attribute) {
	start := len(attrstr.String())
	end := start + len(what)
	attrstr.AppendUnattributed(what)
	for _, a := range attrs {
		attrstr.SetAttribute(a, start, end)
	}
}*/

func (it *ItPlot) drawText(p *ui.AreaDrawParams, text string, x, y float64, size float64, cr, cg, cb, ca float64) {
	attrstr := ui.NewAttributedString(text)
	sz := len(text)
	attrstr.SetAttribute(ui.TextColor{cr, cg, cb, ca}, 0, sz)
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
