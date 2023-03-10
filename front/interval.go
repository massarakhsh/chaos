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
	viewSign int
}

func BuildInterval() *ItInterval {
	it := &ItInterval{}
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
	if dt := data.GetData(it.Sign, 0, 4096); dt != nil {
		it.Load(dt)
		return true
	} else {
		return false
	}
}