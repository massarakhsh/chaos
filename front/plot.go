package front

import (
	"github/massarakhsh/chaos/data"

	"github.com/andlabs/ui"
)

type ItPlot struct{}

func (ItPlot) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	// fill the area with white
	brush := mkSolidBrush(colorWhite, 1.0)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(0, 0, p.AreaWidth, p.AreaHeight)
	path.End()
	p.Context.Fill(path, brush)
	path.Free()

	graphWidth, graphHeight := graphSize(p.AreaWidth, p.AreaHeight)

	sp := &ui.DrawStrokeParams{
		Cap:        ui.DrawLineCapFlat,
		Join:       ui.DrawLineJoinMiter,
		Thickness:  2,
		MiterLimit: ui.DrawDefaultMiterLimit,
	}

	// draw the axes
	brush = mkSolidBrush(colorBlack, 1.0)
	path = ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigure(xoffLeft, yoffTop)
	path.LineTo(xoffLeft, yoffTop+graphHeight)
	path.NewFigure(xoffLeft, yoffTop+graphHeight/2)
	path.LineTo(xoffLeft+graphWidth, yoffTop+graphHeight/2)
	path.End()
	p.Context.Stroke(path, brush, sp)
	path.Free()

	// now transform the coordinate space so (0, 0) is the top-left corner of the graph
	m := ui.DrawNewMatrix()
	m.Translate(xoffLeft, yoffTop)
	p.Context.Transform(m)

	// now get the color for the graph itself and set up the brush
	graphR, graphG, graphB, graphA := colorButton.Color()
	brush.Type = ui.DrawBrushTypeSolid
	brush.R = graphR
	brush.G = graphG
	brush.B = graphB
	// we set brush.A below to different values for the fill and stroke

	// now create the fill for the graph below the graph line
	/*path = constructGraph(graphWidth, graphHeight, true)
	brush.A = graphA / 2
	p.Context.Fill(path, brush)
	path.Free()*/

	// now draw the histogram line
	path = constructGraph(graphWidth, graphHeight, false)
	brush.A = graphA
	p.Context.Stroke(path, brush, sp)
	path.Free()

	// now draw the point being hovered over
	if currentPoint != -1 {
		xs, ys := pointLocations(graphWidth, graphHeight)
		path = ui.DrawNewPath(ui.DrawFillModeWinding)
		path.NewFigureWithArc(
			xs[currentPoint], ys[currentPoint],
			pointRadius,
			0, 6.23, // TODO pi
			false)
		path.End()
		// use the same brush as for the histogram lines
		p.Context.Fill(path, brush)
		path.Free()
	}
}

func inPoint(x, y float64, xtest, ytest float64) bool {
	// TODO switch to using a matrix
	x -= xoffLeft
	y -= yoffTop
	return (x >= xtest-pointRadius) &&
		(x <= xtest+pointRadius) &&
		(y >= ytest-pointRadius) &&
		(y <= ytest+pointRadius)
}

func (ItPlot) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	graphWidth, graphHeight := graphSize(me.AreaWidth, me.AreaHeight)
	xs, ys := pointLocations(graphWidth, graphHeight)

	currentPoint = -1
	for i := 0; i < len(xs); i++ {
		if inPoint(me.X, me.Y, xs[i], ys[i]) {
			currentPoint = i
			break
		}
	}

	// TODO only redraw the relevant area
	histogram.QueueRedrawAll()
}

func (ItPlot) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (ItPlot) DragBroken(a *ui.Area) {
	// do nothing
}

func (ItPlot) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}

func pointLocations(width, height float64) ([]float64, []float64) {
	npoint := data.Serial.Length
	xs := make([]float64, npoint)
	ys := make([]float64, npoint)
	for i := 0; i < npoint; i++ {
		val := data.Serial.Data[i]
		xs[i] = float64(i) * width / float64(npoint)
		ys[i] = (height - val*height) / 2
	}
	return xs, ys
}

func constructGraph(width, height float64, extend bool) *ui.DrawPath {
	xs, ys := pointLocations(width, height)
	path := ui.DrawNewPath(ui.DrawFillModeWinding)

	path.NewFigure(xs[0], ys[0])
	for i := 1; i < len(xs); i++ {
		path.LineTo(xs[i], ys[i])
	}

	if extend {
		path.LineTo(width, height/2)
		path.LineTo(0, height/2)
		path.CloseFigure()
	}

	path.End()
	return path
}
