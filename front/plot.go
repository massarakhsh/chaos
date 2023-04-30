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

	Mouse IfMouse
}

type IfMouse interface {
	RunMouse(nb int, x, y float64, on bool)
}

func (it *ItPlot) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	it.Gate.Lock()
	it.resize(p)
	it.clear(p)
	it.calc(p)
	it.drawFon(p)
	it.drawAxes(p)
	it.drawGraph(p)
	it.drawPens(p)
	it.Gate.Unlock()
}

func (it *ItPlot) clear(p *ui.AreaDrawParams) {
	it.SP = &ui.DrawStrokeParams{
		Cap:        ui.DrawLineCapFlat,
		Join:       ui.DrawLineJoinMiter,
		Thickness:  2,
		MiterLimit: ui.DrawDefaultMiterLimit,
	}
	brush := mkSolidBrush(0xcccccc, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(0, 0, it.Width, it.Height)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()
}

func (it *ItPlot) calc(p *ui.AreaDrawParams) {
	it.X.calibrate(bd, it.Width-2*bd, true)
	it.Y.calibrate(bd, it.Height-2*bd, false)
	it.X.LocZero = it.X.LocDep
	it.Y.LocZero = it.Y.LocDep
}

func (it *ItPlot) drawFon(p *ui.AreaDrawParams) {
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(it.X.LocDep, it.Y.LocDep, it.X.LocSize, it.Y.LocSize)
	path.End()
	brush := mkSolidBrush(0xffffff, 1.0)
	p.Context.Fill(path, brush)
	path.Free()
	if it.CropSign != 0 {
		xf := it.X.ToLoc(it.CropFrom)
		if xf < it.X.LocDep {
			xf = it.X.LocDep
		} else if xf > it.X.LocDep+it.X.LocSize {
			xf = it.X.LocDep + it.X.LocSize
		}
		xt := it.X.ToLoc(it.CropTo)
		if xt < xf {
			xt = xf
		} else if xt > it.X.LocDep+it.X.LocSize {
			xt = it.X.LocDep + it.X.LocSize
		}
		path := ui.DrawNewPath(ui.DrawFillModeWinding)
		path.AddRectangle(xf, it.Y.LocDep, xt-xf, it.Y.LocSize)
		path.End()
		brush := mkSolidBrush(0xccccff, 1.0)
		p.Context.Fill(path, brush)
		path.Free()
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
	if count := len(it.Data); count >= 2 {
		xf := it.X.ToLoc(it.X.Min)
		ixf := int64(xf)
		ymin := it.Y.ToLoc(it.Data[0])
		ymax := ymin
		path.NewFigure(xf, ymin)
		for i := 0; i < count; i++ {
			xt := it.X.ToLoc((it.X.Min*float64(count-1-i) + it.X.Max*float64(i)) / float64(count-1))
			yt := it.Y.ToLoc(it.Data[i])
			if ixt := int64(xt); ixt != ixf {
				path.LineTo(xf, ymin)
				path.LineTo(xf, ymax)
				xf = xt
				ixf = ixt
				ymin = yt
				ymax = ymin
			} else if yt < ymin {
				ymin = yt
			} else if yt > ymax {
				ymax = yt
			}
		}
	}
	path.End()
	p.Context.Stroke(path, brush, it.SP)
	path.Free()
}

func (it *ItPlot) drawPens(p *ui.AreaDrawParams) {
}

func (it *ItPlot) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	if it.Mouse != nil {
		if me.Down > 0 {
			it.Mouse.RunMouse(int(me.Down), me.X, me.Y, true)
		}
		if me.Up > 0 {
			it.Mouse.RunMouse(int(me.Up), me.X, me.Y, false)
		}
	}
	a.QueueRedrawAll()
}

func (it *ItPlot) MouseCrossed(a *ui.Area, left bool) {
}

func (it *ItPlot) DragBroken(a *ui.Area) {
}

func (it *ItPlot) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
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
