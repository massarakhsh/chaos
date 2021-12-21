package front

import (
	"math/rand"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var ()

func MainStart() {
	ui.Main(mainStart)
}

func mainStart() {
	mainwin := ui.NewWindow("ХАОС. Обработка временных серий", 800, 600, true)
	mainwin.SetMargined(true)
	mainwin.OnClosing(func(*ui.Window) bool {
		mainwin.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	mainwin.SetChild(hbox)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	histogram = ui.NewArea(ItPlot{})

	rand.Seed(time.Now().Unix())
	/*for i := 0; i < 10; i++ {
		datapoints[i] = ui.NewSpinbox(0, 100)
		datapoints[i].SetValue(rand.Intn(101))
		datapoints[i].OnChanged(func(*ui.Spinbox) {
			histogram.QueueRedrawAll()
		})
		vbox.Append(datapoints[i], false)
	}*/

	colorButton = ui.NewColorButton()
	// TODO inline these
	brush := mkSolidBrush(colorDodgerBlue, 1.0)
	colorButton.SetColor(brush.R,
		brush.G,
		brush.B,
		brush.A)
	colorButton.OnChanged(func(*ui.ColorButton) {
		histogram.QueueRedrawAll()
	})
	vbox.Append(colorButton, false)

	hbox.Append(histogram, true)

	mainwin.Show()
}

var (
	histogram   *ui.Area
	colorButton *ui.ColorButton

	currentPoint = -1
)

// some metrics
const (
	xoffLeft    = 20 // histogram margins
	yoffTop     = 20
	xoffRight   = 20
	yoffBottom  = 20
	pointRadius = 5
)

// helper to quickly set a brush color
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

// and some colors
// names and values from https://msdn.microsoft.com/en-us/library/windows/desktop/dd370907%28v=vs.85%29.aspx
const (
	colorWhite      = 0xFFFFFF
	colorBlack      = 0x000000
	colorDodgerBlue = 0x1E90FF
)

func graphSize(clientWidth, clientHeight float64) (graphWidth, graphHeight float64) {
	return clientWidth - xoffLeft - xoffRight,
		clientHeight - yoffTop - yoffBottom
}
