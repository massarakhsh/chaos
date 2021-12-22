package front

import (
	"github/massarakhsh/chaos/data"

	"github.com/andlabs/ui"
)

const bd = 10
const radius = 5

type ItSource interface {
	LoadData(sign int) *data.ItSerial
}

type ItPlot struct {
	SP            *ui.DrawStrokeParams
	Width, Height float64
	XZero, YZero  float64
	data          *ItData
}

type ItData struct {
	Sign  int
	Count int
	List  []ItPoint
}

type ItPoint struct {
	ValX, ValY float64
	LocX, LocY float64
}

func (it ItPlot) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	it.resize(p)
	it.clear(p)

	it.SP = &ui.DrawStrokeParams{
		Cap:        ui.DrawLineCapFlat,
		Join:       ui.DrawLineJoinMiter,
		Thickness:  2,
		MiterLimit: ui.DrawDefaultMiterLimit,
	}

	it.drawAxes(p)
	it.drawGraph(p)
	it.drawPens(p)
}

func (it ItPlot) SetSerial(serial *data.ItSerial) {
	if serial == nil || serial.Length <= 1 {
		it.data = nil
	} else {
		length := serial.Length
		data := &ItData{Sign: serial.Sign}
		data.Count = length
		data.List = make([]ItPoint, data.Count)
		for i := 0; i < length; i++ {
			var point ItPoint
			point.ValX = (serial.XMin*float64(length-1-i) + serial.XMax*float64(i)) / float64(length-1)
			point.ValY = serial.Data[i]
			data.List[i] = point
		}
		it.data = data
	}
}

func (it ItPlot) rescan(serial *data.ItSerial) {
	if serial == nil {
		it.data = nil
	} else {
		data := &ItData{}
		data.Sign = serial.Sign
		data.Count = serial.Length
		data.List = make([]ItPoint, data.Count)
		for i := 0; i < data.Count; i++ {
			val := serial.Data[i]
			xs[i] = float64(i) * width / float64(npoint)
			ys[i] = (height - val*height) / 2
		}
		it.data = data
	}
}

func (it *ItPlot) resize(p *ui.AreaDrawParams) {
	it.Width, it.Height = p.AreaWidth-2*bd, p.AreaHeight-2*bd
	m := ui.DrawNewMatrix()
	m.Translate(bd, bd)
	p.Context.Transform(m)
	it.XZero, it.YZero = 0, it.Height/2
}

func (it *ItPlot) clear(p *ui.AreaDrawParams) {
	brush := mkSolidBrush(colorWhite, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(-bd, -bd, it.Width+2*bd, it.Height+2*bd)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()
}

func (it *ItPlot) drawAxes(p *ui.AreaDrawParams) {
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

func (it *ItPlot) drawGraph(p *ui.AreaDrawParams) {
	brush := mkSolidBrush(colorDodgerBlue, 0.5)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)

	xs, ys := pointLocations(it.Width, it.Height)
	if len(xs) > 0 {
		path.NewFigure(xs[0], ys[0])
		for i := 1; i < len(xs); i++ {
			path.LineTo(xs[i], ys[i])
		}
	}
	path.End()
	p.Context.Stroke(path, brush, it.SP)
	path.Free()

}

func (it *ItPlot) drawPens(p *ui.AreaDrawParams) {
	if data.Ser.Current >= 0 {
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
	}
}

func inPoint(x, y float64, xtest, ytest float64) bool {
	// TODO switch to using a matrix
	return (x >= xtest-radius) &&
		(x <= xtest+radius) &&
		(y >= ytest-radius) &&
		(y <= ytest+radius)
}

func (it ItPlot) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	xs, ys := pointLocations(it.Width, it.Height)

	data.Ser.Current = -1
	for i := 0; i < len(xs); i++ {
		if inPoint(me.X, me.Y, xs[i], ys[i]) {
			data.Ser.Current = i
			break
		}
	}

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

func pointLocations(width, height float64) ([]float64, []float64) {
	npoint := data.Ser.Length
	xs := make([]float64, npoint)
	ys := make([]float64, npoint)
	for i := 0; i < npoint; i++ {
		val := data.Ser.Data[i]
		xs[i] = float64(i) * width / float64(npoint)
		ys[i] = (height - val*height) / 2
	}
	return xs, ys
}
