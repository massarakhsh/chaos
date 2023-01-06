package front

import (
	"math/rand"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/massarakhsh/chaos/data"
)

type ItWindow struct {
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
	if combo := ui.NewCombobox(); combo != nil {
		combo.Append("Останов")
		combo.Append("Модель")
		combo.Append("COM-порт")
		combo.OnSelected(func(c *ui.Combobox) {
			if sel := c.Selected(); sel == 1 {
				data.SetSource(data.SOURCE_MODEL)
			} else if sel == 2 {
				data.SetSource(data.SOURCE_SERIAL)
			} else {
				data.SetSource(data.SOURCE_NO)
			}
		})
		box.Append(combo, false)
		combo.SetSelected(2)
	}
	if combo := ui.NewCombobox(); combo != nil {
		combo.Append("Пошагово")
		combo.Append("1 Гц")
		combo.Append("5 Гц")
		combo.OnSelected(func(c *ui.Combobox) {
			if sel := c.Selected(); sel == 1 {
				DuraUpdate = 1000
			} else if sel == 2 {
				DuraUpdate = 200
			} else {
				DuraUpdate = 0
			}
		})
		box.Append(combo, false)
		combo.SetSelected(1)
	}
	if button := ui.NewButton("Обновить"); button != nil {
		button.OnClicked(func(b *ui.Button) {
			NeedUpdate = 0x3
		})
		box.Append(button, false)
	}
	if button := ui.NewButton("Сброс"); button != nil {
		button.OnClicked(func(b *ui.Button) {
			data.GenReset()
		})
		box.Append(button, false)
	}
	return box
}

func (it *ItWindow) buildInfoBox() *ui.Box {
	box := ui.NewVerticalBox()
	box.SetPadded(true)
	if elm := it.buildPlotBox(); elm != nil {
		box.Append(elm, true)
	}
	if elm := it.buildSpectrBox(); elm != nil {
		box.Append(elm, true)
	}
	return box
}

func (it *ItWindow) buildPlotBox() *ui.Box {
	box := ui.NewHorizontalBox()
	box.SetPadded(true)
	graph := ui.NewArea(BuildGraphic())
	box.Append(graph, true)
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
					if MainGraphic.Probe() {
						graph.QueueRedrawAll()
					}
				}
			}
		}
	}()
	return box
}

func (it *ItWindow) buildSpectrBox() *ui.Box {
	box := ui.NewVerticalBox()
	box.SetPadded(true)
	spectr := ui.NewArea(BuildSpectr())
	box.Append(spectr, true)
	go func() {
		nexttime := time.Now()
		for {
			time.Sleep(time.Millisecond * 10)
			if time.Now().After(nexttime) {
				if DuraUpdate > 0 {
					nexttime = time.Now().Add(time.Millisecond * time.Duration(DuraUpdate))
				}
				if DuraUpdate > 0 || (NeedUpdate&0x2) != 0 {
					NeedUpdate &= 0xfd
					if MainSpectr.Probe() {
						spectr.QueueRedrawAll()
					}
				}
			}
		}
	}()
	return box
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
