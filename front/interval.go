package front

import (
	"time"

	"github.com/andlabs/ui"
	"github.com/massarakhsh/chaos/data"
	"github.com/massarakhsh/chaos/pkg/zone"
)

type ItInterval struct {
	zone.ItZone
	ItPlot

	Source   *ItPanel
	ViewSign int
}

func BuildInterval(source *ItPanel) *ItInterval {
	it := &ItInterval{Source: source}
	it.area = ui.NewArea(it)
	it.BindControl(it, it.area)
	it.Panel.Loader = it
	it.Panel.Width, it.Panel.Height = 1024, 512
	it.BindRefresh(it)
	return it
}

func (it *ItInterval) Refresh() {
	if time.Since(it.lastUpdate) >= time.Second*1 {
		if it.Probe() {
			it.area.QueueRedrawAll()
		}
		it.lastUpdate = time.Now()
	}
}

func (it *ItInterval) Probe() bool {
	if it.Panel.Sign != it.Source.Sign || it.ViewSign != it.Source.CropSign {
		if dt := it.getData(); dt != nil {
			it.Panel.Load(dt)
			return true
		}
	}
	return false
}

func (it *ItInterval) getData() *data.ItData {
	it.ViewSign = it.Source.CropSign
	return it.Source.GetData()
}
