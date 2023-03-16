package front

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/massarakhsh/chaos/data"
	"github.com/massarakhsh/chaos/pkg/zone"
)

type ItControl struct {
	zone.ItZone
}

func buildControl() *ItControl {
	it := &ItControl{}
	it.BindVerticalBox(it)
	it.addMode()
	return it
}

func (it *ItControl) addMode() {
	level := it.GetVolume()
	if combo := ui.NewCombobox(); combo != nil {
		combo.Append("Сброс")
		combo.Append("Чтение COM-порта")
		combo.Append("Моделирование")
		combo.Append("Прочитать файл")
		combo.Append("Анализ")
		combo.SetSelected(0)
		combo.OnSelected(func(c *ui.Combobox) {
			it.setMode(level, c.Selected())
		})
		it.PushControl(combo)
		it.setMode(level, 0)
	}
}

func (it *ItControl) setMode(level int, sel int) {
	it.PopControls(level)
	if sel == 1 {
		data.SetSource(data.SOURCE_SERIAL)
	} else if sel == 2 {
		data.SetSource(data.SOURCE_MODEL)
		it.addTemp()
	} else if sel == 3 {
		data.SetSource(data.SOURCE_FILE)
		it.addReadFile()
	} else if sel == 4 {
		data.SetSource(data.SOURCE_ANALIZE)
	} else {
		data.SetSource(data.SOURCE_RESET)
	}
}

func (it *ItControl) setModeAnalize() {
	it.setMode(0, 4)
}

func (it *ItControl) addTemp() {
	level := it.GetVolume()
	if combo := ui.NewCombobox(); combo != nil {
		combo.Append("Пошагово")
		combo.Append("1 Гц")
		combo.Append("5 Гц")
		combo.SetSelected(1)
		combo.OnSelected(func(c *ui.Combobox) {
			it.setTemp(level, c.Selected())
		})
		it.PushControl(combo)
		it.setTemp(level, 1)
	}
}

func (it *ItControl) setTemp(level int, sel int) {
	it.PopControls(level)
	if sel == 1 {
		DuraUpdate = 1000
	} else if sel == 2 {
		DuraUpdate = 200
	} else {
		DuraUpdate = 0
	}
}

func (it *ItControl) addReadFile() {
	//level := it.GetVolume()
	if button := ui.NewButton("Прочитать"); button != nil {
		button.OnClicked(func(b *ui.Button) {
			data.ReadFromFile()
			it.setModeAnalize()
		})
		it.PushControl(button)
	}
}
