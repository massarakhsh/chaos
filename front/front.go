package front

import (
	"math/rand"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/massarakhsh/chaos/pkg/zone"
)

type ItFront struct {
	zone.ItZone
}

var IsTerminating bool
var IsAutoView bool
var IsSignView bool

func MainStart() {
	it := &ItFront{}
	go it.run()
	ui.Main(it.mainStart)
}

func (it *ItFront) mainStart() {
	rand.Seed(time.Now().Unix())
	if mainWindow := ui.NewWindow("ХАОС. Обработка временных серий", 1280, 800, true); mainWindow != nil {
		mainWindow.SetMargined(true)
		mainWindow.OnClosing(func(*ui.Window) bool {
			mainWindow.Destroy()
			ui.Quit()
			IsTerminating = true
			return false
		})
		it.BindHorizontalBox(it)
		mainWindow.SetChild(it.GetControl())
		it.buildMain()

		mainWindow.Show()
	}
}

func (it *ItFront) run() {
	time.Sleep(time.Second * 1)
	for !IsTerminating {
		time.Sleep(time.Millisecond * 1000)
		it.Step()
	}
}

func (it *ItFront) buildMain() {
	if left := zone.BuildVerticalBox(nil); left != nil {
		it.Append(left, false)
		if child := buildControl(); child != nil {
			left.Append(child, true)
		}
		if child := buildFile(); child != nil {
			left.Append(child, false)
		}
	}
	if right := zone.BuildVerticalBox(nil); right != nil {
		it.Append(right, true)
		if child := buildInfo(); child != nil {
			right.Append(child, true)
		}
	}
}

func SignalRedraw() {
	time.AfterFunc(time.Millisecond*1500, func() { IsSignView = true })
}
