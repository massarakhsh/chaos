package front

import (
	"math/rand"
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
		dt = data.GetData(it.Panel.Sign, 0, 65536*4)
	} else if it.Panel.Sign != ViewSign {
		if dt = data.GetData(it.Panel.Sign, 0, 65536*4); dt != nil {
			dt.Sign = ViewSign
		}
	}
	if dt != nil {
		it.Panel.Load(dt)
		return true
	}
	return false
}

func (it *ItSignal) RunMouse(nb int, x, y float64, on bool) {
	val := it.Panel.X.ToVal(x)
	if nb == 1 && on {
		it.Panel.CropSign = rand.Int()
		it.Panel.CropFrom = val
		if it.Panel.CropTo < it.Panel.CropFrom {
			it.Panel.CropTo = it.Panel.X.Max
		}
	} else if nb == 3 && on {
		it.Panel.CropSign = rand.Int()
		it.Panel.CropTo = val
		if it.Panel.CropFrom > it.Panel.CropTo {
			it.Panel.CropFrom = it.Panel.X.Min
		}
	} else if nb == 2 && on {
		it.Panel.CropSign = 0
	}
}
