package front

import (
	"math/rand"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

type ItFront struct {
	MainWindow   *ui.Window
	MainBox      *ui.Box
	ControlBox   *ui.Box
	ListControls []ItControl

	InfoBox   *ui.Box
	InfoArea  *ui.Area
	RowBox    *ui.Box
	ListAreas []*ui.Area
}

var DuraUpdate = 1000
var NeedUpdate = 0x3

func MainStart() {
	it := &ItFront{}
	ui.Main(it.mainStart)
}

func (it *ItFront) mainStart() {
	rand.Seed(time.Now().Unix())
	it.MainWindow = ui.NewWindow("ХАОС. Обработка временных серий", 800, 600, true)
	it.MainWindow.SetMargined(true)
	it.MainWindow.OnClosing(func(*ui.Window) bool {
		it.MainWindow.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		it.MainWindow.Destroy()
		return true
	})

	it.buildMainBox()
	it.addControlMode()

	it.MainWindow.Show()
}

func (it *ItFront) buildMainBox() {
	it.MainBox = ui.NewHorizontalBox()
	it.MainBox.SetPadded(true)
	it.MainWindow.SetChild(it.MainBox)

	it.ControlBox = ui.NewVerticalBox()
	it.ControlBox.SetPadded(true)
	it.MainBox.Append(it.ControlBox, false)

	it.InfoBox = ui.NewVerticalBox()
	it.InfoBox.SetPadded(true)
	it.MainBox.Append(it.InfoBox, true)
}
