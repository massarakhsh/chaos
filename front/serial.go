package front

import (
	"github/massarakhsh/chaos/data"

	"github.com/andlabs/ui"
)

const bd = 10
const radius = 5

type ItLoad interface {
	Probe() bool
}

type ItSerial struct {
	Sign          int
	SP            *ui.DrawStrokeParams
	Width, Height float64
	Count         int
	List          []ItPoint
	XMin, XMax    float64
	YMin, YMax    float64
	XZero, YZero  float64

	Loader ItLoad
}

type ItPoint struct {
	XVal, YVal float64
	XLoc, YLoc float64
}

func (it *ItSerial) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	it.Loader.Probe()
	it.resize(p)
	it.clear(p)
	it.calc(p)
	it.drawAxes(p)
	it.drawGraph(p)
	it.drawPens(p)
}

func (it *ItSerial) resize(p *ui.AreaDrawParams) {
	it.Width, it.Height = p.AreaWidth-2*bd, p.AreaHeight-2*bd
	m := ui.DrawNewMatrix()
	m.Translate(bd, bd)
	p.Context.Transform(m)
}

func (it *ItSerial) clear(p *ui.AreaDrawParams) {
	it.SP = &ui.DrawStrokeParams{
		Cap:        ui.DrawLineCapFlat,
		Join:       ui.DrawLineJoinMiter,
		Thickness:  2,
		MiterLimit: ui.DrawDefaultMiterLimit,
	}
	brush := mkSolidBrush(colorWhite, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(-bd, -bd, it.Width+2*bd, it.Height+2*bd)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()
}

func (it *ItSerial) calc(p *ui.AreaDrawParams) {
	it.XZero, it.YZero = 0, it.Height
	if length := it.Count; length >= 2 {
		x0 := it.XMin * it.Width / (it.XMin - it.XMax)
		if x0 < 0 {
			x0 = 0
		} else if x0 > it.Width {
			x0 = it.Width
		}
		it.XZero = x0
		y0 := it.YMax * it.Height / (it.YMax - it.YMin)
		if y0 < 0 {
			y0 = 0
		} else if y0 > it.Height {
			y0 = it.Height
		}
		it.YZero = y0
		for i := 0; i < length; i++ {
			point := &it.List[i]
			point.XLoc = (0*(it.XMax-point.XVal) + it.Width*(point.XVal-it.XMin)) / (it.XMax - it.XMin)
			point.YLoc = (0*(point.YVal-it.YMin) + it.Height*(it.YMax-point.YVal)) / (it.YMax - it.YMin)
		}
	}
}

func (it *ItSerial) drawAxes(p *ui.AreaDrawParams) {
	brush := mkSolidBrush(colorBlack, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigure(it.XZero, 0)
	path.LineTo(it.XZero, it.Height)
	path.NewFigure(0, it.YZero)
	path.LineTo(it.Width, it.YZero)
	path.End()
	p.Context.Stroke(path, brush, it.SP)
	path.Free()
}

func (it *ItSerial) drawGraph(p *ui.AreaDrawParams) {
	brush := mkSolidBrush(colorDodgerBlue, 0.5)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigure(it.List[0].XLoc, it.List[0].YLoc)
	for i := 1; i < it.Count; i++ {
		path.LineTo(it.List[i].XLoc, it.List[i].YLoc)
	}
	path.End()
	p.Context.Stroke(path, brush, it.SP)
	path.Free()
}

func (it *ItSerial) drawPens(p *ui.AreaDrawParams) {
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

func (it *ItSerial) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
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

func (it *ItSerial) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (it *ItSerial) DragBroken(a *ui.Area) {
	// do nothing
}

func (it *ItSerial) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}

func (it *ItSerial) Load(serial *data.ItData) {
	if serial == nil || serial.Length < 2 {
		it.Count = 0
		it.List = []ItPoint{}
	} else {
		length := serial.Length
		it.Sign = serial.Sign
		it.Count = length
		it.List = make([]ItPoint, length)
		it.XMin = serial.XMin
		it.XMax = serial.XMax
		if it.XMin >= it.XMax {
			it.XMin -= 0.1
			it.XMax += 0.1
		}
		for i := 0; i < length; i++ {
			point := &it.List[i]
			point.XVal = (serial.XMin*float64(length-1-i) + serial.XMax*float64(i)) / float64(length-1)
			val := serial.Data[i]
			point.YVal = val
			if i == 0 || val < it.YMin {
				it.YMin = val
			}
			if i == 0 || val > it.YMax {
				it.YMax = val
			}
		}
		if it.YMin >= 0 {
			it.YMin = -it.YMax
		} else if it.YMax <= 0 {
			it.YMax = -it.YMin
		} else if it.YMin > -it.YMax {
			it.YMin = -it.YMax
		} else if it.YMax < -it.YMin {
			it.YMax = -it.YMin
		}
		if it.YMin >= it.YMax {
			it.YMin -= 0.1
			it.YMax += 0.1
		}
	}
}
