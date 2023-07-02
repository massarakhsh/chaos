package front

import (
	"fmt"
	"time"

	"github.com/andlabs/ui"
)

const bd = 10

//const radius = 5

type ItPlot struct {
	Panel      ItPanel
	SP         *ui.DrawStrokeParams
	area       *ui.Area
	lastUpdate time.Time

	Mouse IfMouse
}

type IfMouse interface {
	RunMouse(nb int, x, y float64, on bool)
}

func (it *ItPlot) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	it.Panel.Gate.Lock()
	it.Panel.resize(p)
	it.clear(p)
	it.calc(p)
	it.drawFon(p)
	it.drawAxes(p)
	it.drawGraph(p)
	it.drawPens(p)
	it.Panel.Gate.Unlock()
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
	path.AddRectangle(0, 0, it.Panel.Width, it.Panel.Height)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()
}

func (it *ItPlot) calc(p *ui.AreaDrawParams) {
	it.Panel.X.calibrate(bd, it.Panel.Width-2*bd, true)
	it.Panel.Y.calibrate(bd, it.Panel.Height-2*bd, false)
	it.Panel.X.LocZero = it.Panel.X.LocDep
	it.Panel.Y.LocZero = it.Panel.Y.LocDep
}

func (it *ItPlot) drawFon(p *ui.AreaDrawParams) {
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(it.Panel.X.LocDep, it.Panel.Y.LocDep, it.Panel.X.LocSize, it.Panel.Y.LocSize)
	path.End()
	brush := mkSolidBrush(0xffffff, 1.0)
	p.Context.Fill(path, brush)
	path.Free()
	if it.Panel.CropSign != 0 {
		xf := it.Panel.X.ToLoc(it.Panel.CropFrom)
		if xf < it.Panel.X.LocDep {
			xf = it.Panel.X.LocDep
		} else if xf > it.Panel.X.LocDep+it.Panel.X.LocSize {
			xf = it.Panel.X.LocDep + it.Panel.X.LocSize
		}
		xt := it.Panel.X.ToLoc(it.Panel.CropTo)
		if xt < xf {
			xt = xf
		} else if xt > it.Panel.X.LocDep+it.Panel.X.LocSize {
			xt = it.Panel.X.LocDep + it.Panel.X.LocSize
		}
		path := ui.DrawNewPath(ui.DrawFillModeWinding)
		path.AddRectangle(xf, it.Panel.Y.LocDep, xt-xf, it.Panel.Y.LocSize)
		path.End()
		brush := mkSolidBrush(0xccccff, 1.0)
		p.Context.Fill(path, brush)
		path.Free()
	}
}

func (it *ItPlot) drawAxes(p *ui.AreaDrawParams) {
	if path := ui.DrawNewPath(ui.DrawFillModeWinding); path != nil {
		for x := it.Panel.X.First; x < it.Panel.X.Max; x += it.Panel.X.Step {
			if xt := it.Panel.X.ToLoc(x); xt >= 0 && xt < it.Panel.Width {
				path.NewFigure(it.Panel.X.LocDep+xt, it.Panel.X.LocDep+it.Panel.Y.LocSize)
				path.LineTo(it.Panel.X.LocDep+xt, it.Panel.X.LocDep)
				it.drawText(p, fmt.Sprintf(it.Panel.X.Format, x), it.Panel.X.LocDep+xt-16, it.Panel.Y.LocDep+it.Panel.Y.LocSize-10, 10, 0.99, 0, 0, 0.99)
			}
		}
		path.End()
		brush := mkSolidBrush(0x808080, 0.5)
		p.Context.Stroke(path, brush, it.SP)
		path.Free()
	}
	if path := ui.DrawNewPath(ui.DrawFillModeWinding); path != nil {
		path.NewFigure(it.Panel.X.LocZero, it.Panel.Y.LocDep+it.Panel.Y.LocSize)
		path.LineTo(it.Panel.X.LocZero, it.Panel.Y.LocDep)
		path.NewFigure(it.Panel.X.LocDep, bd+it.Panel.Y.LocSize-it.Panel.Y.LocZero)
		path.LineTo(it.Panel.X.LocDep+it.Panel.X.LocSize, bd+it.Panel.Y.LocSize-it.Panel.Y.LocZero)
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
		Width:       it.Panel.Width - x,
		Align:       ui.DrawTextAlignLeft,
	}
	tl := ui.DrawNewTextLayout(&lp)
	p.Context.Text(tl, x, y)
	tl.Free()
}

func (it *ItPlot) drawGraph(p *ui.AreaDrawParams) {
	brush := mkSolidBrush(0x5599ff, 0.5)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	if count := len(it.Panel.Data); count >= 2 {
		xf := it.Panel.X.ToLoc(it.Panel.X.Min)
		ixf := int64(xf)
		ymin := it.Panel.Y.ToLoc(it.Panel.Data[0])
		ymax := ymin
		path.NewFigure(xf, ymin)
		for i := 0; i < count; i++ {
			xt := it.Panel.X.ToLoc((it.Panel.X.Min*float64(count-1-i) + it.Panel.X.Max*float64(i)) / float64(count-1))
			yt := it.Panel.Y.ToLoc(it.Panel.Data[i])
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
