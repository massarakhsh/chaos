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
	it.Loader = it
	it.IsZeroCenter = true
	it.Width, it.Height = 1024, 512
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
	if it.Sign != it.Signal.Sign || it.ViewSign != it.Signal.CropSign {
		if dt := it.getData(); dt != nil {
			it.Load(dt)
			return true
		}
	}
	return false
}

func (it *ItInterval) getData() *data.ItData {
	it.ViewSign = it.Signal.CropSign
	return it.Signal.GetData()
}
