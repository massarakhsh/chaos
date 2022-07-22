package front

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

const (
	colorWhite      = 0xFFFFFF
	colorBlack      = 0x000000
	colorDodgerBlue = 0x1E90FF
)

type ItWindow struct {
	Window     *ui.Window
	MainBox    *ui.Box
	ControlBox *ui.Box
	InfoBox    *ui.Box
	DownBox    *ui.Box
	Graph      *ui.Area
	Org        *ui.Area
	Spectr     *ui.Area
	ButRes     *ui.Button
	ButColRes  *ui.ColorButton
}

func MainStart() {
	it := &ItWindow{}
	ui.Main(it.mainStart)
}

func (it *ItWindow) mainStart() {
	mainwin := ui.NewWindow("ХАОС. Обработка временных серий", 800, 600, true)
	it.Window = mainwin

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

	it.MainBox = ui.NewHorizontalBox()
	it.MainBox.SetPadded(true)
	mainwin.SetChild(it.MainBox)

	it.ControlBox = ui.NewVerticalBox()
	it.ControlBox.SetPadded(true)
	it.MainBox.Append(it.ControlBox, false)

	it.InfoBox = ui.NewVerticalBox()
	it.InfoBox.SetPadded(true)
	it.MainBox.Append(it.InfoBox, true)

	it.Graph = ui.NewArea(BuildGraphic())
	it.InfoBox.Append(it.Graph, true)

	it.DownBox = ui.NewHorizontalBox()
	it.DownBox.SetPadded(true)
	it.InfoBox.Append(it.DownBox, true)

	it.ButColRes = ui.NewColorButton()
	brush := mkSolidBrush(colorDodgerBlue, 1.0)
	it.ButColRes.SetColor(brush.R,
		brush.G,
		brush.B,
		brush.A)
	it.ButColRes.OnChanged(func(*ui.ColorButton) {
		fmt.Println("Color changed")
	})
	it.DownBox.Append(it.ButColRes, false)

	it.Spectr = ui.NewArea(BuildSpectr())
	it.DownBox.Append(it.Spectr, true)

	it.Org = ui.NewArea(BuildOrg())
	it.DownBox.Append(it.Org, true)

	rand.Seed(time.Now().Unix())
	/*for i := 0; i < 10; i++ {
		datapoints[i] = ui.NewSpinbox(0, 100)
		datapoints[i].SetValue(rand.Intn(101))
		datapoints[i].OnChanged(func(*ui.Spinbox) {
			histogram.QueueRedrawAll()
		})
		vbox.Append(datapoints[i], false)
	}*/

	it.mainMonitor()
	it.Window.Show()
}

func (it *ItWindow) mainMonitor() {
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			if it.Graph != nil {
				if MainGraphic.Probe() {
					it.Graph.QueueRedrawAll()
				}
			}
			if it.Spectr != nil {
				if MainSpectr.Probe() {
					it.Spectr.QueueRedrawAll()
				}
			}
		}
	}()
}

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
