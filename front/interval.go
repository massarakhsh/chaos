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

	Signal   *ItSignal
	ViewSign int
}

func BuildInterval(signal *ItSignal) *ItInterval {
	it := &ItInterval{Signal: signal}
	it.area = ui.NewArea(it)
	it.BindControl(it, it.area)
	it.Panel.Loader = it
	it.Panel.IsZeroCenter = true
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
	if it.Panel.Sign != it.Signal.Panel.Sign || it.ViewSign != it.Signal.Panel.CropSign {
		if dt := it.getData(); dt != nil {
			it.Panel.Load(dt)
			return true
		}
	}
	return false
}

func (it *ItInterval) getData() *data.ItData {
	it.ViewSign = it.Signal.Panel.CropSign
	return it.Signal.Panel.GetData()
}
