package front

import (
	"math/rand"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

type ItWindow struct {
	Controls *ItControls
}

var DuraUpdate = 1000
var NeedUpdate = 0x3

func MainStart() {
	it := &ItWindow{}
	ui.Main(it.mainStart)
}

func (it *ItWindow) mainStart() {
	rand.Seed(time.Now().Unix())
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

	if mainBox := it.buildMainBox(); mainBox != nil {
		mainwin.SetChild(mainBox)
	}

	mainwin.Show()
}

func (it *ItWindow) buildMainBox() *ui.Box {
	box := ui.NewHorizontalBox()
	box.SetPadded(true)
	if controlBox := it.buildControlBox(); controlBox != nil {
		box.Append(controlBox, false)
	}
	if infoBox := it.buildInfoBox(); infoBox != nil {
		box.Append(infoBox, true)
	}
	return box
}

func (it *ItWindow) buildControlBox() *ui.Box {
	box := ui.NewVerticalBox()
	box.SetPadded(true)
	it.Controls = BuildControls(box)
	return box
}

func (it *ItWindow) buildInfoBox() *ui.Box {
	box := ui.NewVerticalBox()
	box.SetPadded(true)
	if elm := it.buildPlotArea(BuildGraphic()); elm != nil {
		box.Append(elm, true)
	}
	if row := ui.NewHorizontalBox(); row != nil {
		box.Append(row, true)
		if elm := it.buildPlotArea(BuildSignal()); elm != nil {
			row.Append(elm, true)
		}
		if elm := it.buildPlotArea(BuildSpectr()); elm != nil {
			row.Append(elm, true)
		}
	}
	return box
}

func (it *ItWindow) buildPlotArea(plot InfPlot) *ui.Area {
	area := ui.NewArea(plot)
	go func() {
		nexttime := time.Now()
		for {
			time.Sleep(time.Millisecond * 10)
			if time.Now().After(nexttime) {
				if DuraUpdate > 0 {
					nexttime = time.Now().Add(time.Millisecond * time.Duration(DuraUpdate))
				}
				if DuraUpdate > 0 || (NeedUpdate&0x1) != 0 {
					NeedUpdate &= 0xfe
					if plot.Probe() {
						area.QueueRedrawAll()
					}
				}
			}
		}
	}()
	return area
}
