package front

import (
	"fmt"
	"time"

	"github.com/andlabs/ui"
	"github.com/massarakhsh/chaos/data"
	"github.com/massarakhsh/chaos/pkg/zone"
)

type ItFile struct {
	zone.ItZone
	counter    *ui.Label
	lastUpdate time.Time
}

func buildFile() *ItFile {
	it := &ItFile{}
	it.BindVerticalBox(it)
	it.addControls()
	it.BindRefresh(it)
	return it
}

func (it *ItFile) addControls() {
	if button := ui.NewButton("Обновить"); button != nil {
		button.OnClicked(func(b *ui.Button) {
			ViewSign++
		})
		it.PushControl(button)
	}
	if button := ui.NewButton("Сброс"); button != nil {
		button.OnClicked(func(b *ui.Button) {
			data.GenReset()
		})
		it.PushControl(button)
	}
	if it.counter = ui.NewLabel(""); it.counter != nil {
		it.PushControl(it.counter)
	}
	if button := ui.NewButton("Записать"); button != nil {
		button.OnClicked(func(b *ui.Button) {
			data.SaveToFile()
		})
		it.PushControl(button)
	}
}

func (it *ItFile) Refresh() {
	if it.counter.Enabled() && time.Since(it.lastUpdate) >= time.Second*1 {
		text := fmt.Sprintf("%d точек", data.GetDataSize())
		it.counter.SetText(text)
		it.lastUpdate = time.Now()
	}
}
