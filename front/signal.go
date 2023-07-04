package front

import (
	"time"

	"github.com/andlabs/ui"
	"github.com/massarakhsh/chaos/data"
	"github.com/massarakhsh/chaos/pkg/zone"
)

type ItSignal struct {
	zone.ItZone
	ItPlot
}

func BuildSignal() *ItSignal {
	it := &ItSignal{}
	it.area = ui.NewArea(it)
	it.Mouse = it
	it.BindControl(it, it.area)
	it.Panel.Loader = it
	it.Panel.IsZeroCenter = true
	it.Panel.Width, it.Panel.Height = 1024, 512
	it.BindRefresh(it)
	return it
}

func (it *ItSignal) Refresh() {
	if time.Since(it.lastUpdate) >= time.Second*1 {
		if it.Probe() {
			it.area.QueueRedrawAll()
		}
		it.lastUpdate = time.Now()
	}
}

func (it *ItSignal) Probe() bool {
	var dt *data.ItData
	if IsAutoView {
		dt = data.GetData(it.Panel.SignIn, 0, 65536*4)
	} else if IsSignView {
		IsSignView = false
		dt = data.GetData("", 0, 65536*4)
	}
	if dt != nil {
		it.Panel.Load(dt)
		return true
	}
	return false
}
