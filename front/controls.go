package front

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/massarakhsh/chaos/data"
)

type ItControls struct {
	Box      *ui.Box
	Controls []ItControl

	Mode int
}

type ItControl struct {
	Level int
	Con   *ui.ControlBase
}

func BuildControls(box *ui.Box) *ItControls {
	it := &ItControls{Box: box}
	it.addControlMode()
	return it
}

func (it *ItControls) addControlMode() {
	level := len(it.Controls)
	if combo := ui.NewCombobox(); combo != nil {
		combo.Append("Останов")
		combo.Append("Модель")
		combo.Append("COM-порт")
		combo.OnSelected(func(c *ui.Combobox) {
			it.setControlMode(level, c, c.Selected())
		})
		combo.SetSelected(0)
		it.pushControls(&combo.ControlBase)
		it.setControlMode(level, combo, 0)
	}
}

func (it *ItControls) setControlMode(level int, c *ui.Combobox, sel int) {
	it.popControls(level)
	if sel == 1 {
		it.Mode = data.SOURCE_MODEL
		data.SetSource(data.SOURCE_MODEL)
	} else if sel == 2 {
		it.Mode = data.SOURCE_SERIAL
		data.SetSource(data.SOURCE_SERIAL)
	} else {
		it.Mode = data.SOURCE_NO
		data.SetSource(data.SOURCE_NO)
	}
	it.addControlTemp()
}

func (it *ItControls) addControlTemp() {
	level := len(it.Controls)
	if it.Mode == data.SOURCE_MODEL {
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
	}
}

func (it *ItControls) setControlTemp(level int, c *ui.Combobox, sel int) {
	it.popControls(level)
	if sel == 1 {
		DuraUpdate = 1000
	} else if sel == 2 {
		DuraUpdate = 200
	} else {
		DuraUpdate = 0
	}
	/*if button := ui.NewButton("Обновить"); button != nil {
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
	}*/
}

func (it *ItControls) pushControls(con *ui.ControlBase) {
	ctr := ItControl{Level: len(it.Controls), Con: con}
	it.Controls = append(it.Controls, ctr)
	it.Box.Append(con, false)
}

func (it *ItControls) popControls(last int) {
	for pos := len(it.Controls) - 1; pos > last; pos-- {
		it.Box.Delete(pos)
	}
	it.Controls = it.Controls[:last+1]
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
