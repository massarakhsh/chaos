package front

import (
	"fmt"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/massarakhsh/chaos/data"
	"github.com/massarakhsh/chaos/pkg/lit"
)

type ItControl struct {
	Con  *ui.ControlBase
	Sign int64
}

func (it *ItFront) addControlMode() {
	level := len(it.ListControls)
	if combo := ui.NewCombobox(); combo != nil {
		combo.Append("Сброс")
		combo.Append("Чтение COM-порта")
		combo.Append("Моделирование")
		combo.Append("Анализ")
		combo.OnSelected(func(c *ui.Combobox) {
			it.setControlMode(level, c.Selected())
		})
		combo.SetSelected(0)
		it.pushControls(&combo.ControlBase)
	}
}

func (it *ItFront) setControlMode(level int, sel int) {
	it.popControls(level)
	if sel == 1 {
		data.SetSource(data.SOURCE_SERIAL)
	} else if sel == 2 {
		data.SetSource(data.SOURCE_MODEL)
	} else if sel == 3 {
		data.SetSource(data.SOURCE_ANALIZE)
	} else {
		data.SetSource(data.SOURCE_RESET)
	}
	it.addControlTemp()
}

func (it *ItFront) addControlTemp() {
	level := len(it.ListControls)
	if data.GetSource() == data.SOURCE_MODEL {
		if combo := ui.NewCombobox(); combo != nil {
			combo.Append("Пошагово")
			combo.Append("1 Гц")
			combo.Append("5 Гц")
			combo.OnSelected(func(c *ui.Combobox) {
				it.setControlTemp(level, c, c.Selected())
			})
			combo.SetSelected(1)
			it.pushControls(&combo.ControlBase)
			it.setControlTemp(level, combo, 1)
		}
	} else {
		it.addControlWrite()
	}
}

func (it *ItFront) setControlTemp(level int, c *ui.Combobox, sel int) {
	it.popControls(level)
	if sel == 1 {
		DuraUpdate = 1000
	} else if sel == 2 {
		DuraUpdate = 200
	} else {
		DuraUpdate = 0
	}
	it.addControlWrite()
}

func (it *ItFront) addControlWrite() {
	//level := len(it.ListControls)
	if DuraUpdate == 0 {
		if button := ui.NewButton("Обновить"); button != nil {
			button.OnClicked(func(b *ui.Button) {
				NeedUpdate = 0x3
			})
			it.pushControls(&button.ControlBase)
		}
	}
	if button := ui.NewButton("Сброс"); button != nil {
		button.OnClicked(func(b *ui.Button) {
			data.GenReset()
		})
		it.pushControls(&button.ControlBase)
	}
	text := fmt.Sprintf("%d точек", data.GetDataSize())
	if label := ui.NewLabel(text); label != nil {
		go func() {
			var timer lit.ItTimer
			for label.Enabled() {
				if timer.Probe(time.Millisecond*10, time.Second) {
					text := fmt.Sprintf("%d точек", data.GetDataSize())
					label.SetText(text)
				}
			}
		}()
		it.pushControls(&label.ControlBase)
	}
	if button := ui.NewButton("Записать"); button != nil {
		button.OnClicked(func(b *ui.Button) {
			data.SaveToFile()
		})
		it.pushControls(&button.ControlBase)
	}
	it.setInfoMode()
}

func (it *ItFront) pushControls(con *ui.ControlBase) {
	ctr := ItControl{Con: con}
	it.ListControls = append(it.ListControls, ctr)
	it.ControlBox.Append(con, false)
}

func (it *ItFront) popControls(last int) {
	it.purgeInfo()
	for pos := len(it.ListControls) - 1; pos > last; pos-- {
		it.ControlBox.Delete(pos)
	}
	it.ListControls = it.ListControls[:last+1]
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
