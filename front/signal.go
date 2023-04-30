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
	viewSign int
}

func BuildSignal() *ItSignal {
	it := &ItSignal{}
	it.area = ui.NewArea(it)
	it.Mouse = it
	it.BindControl(it, it.area)
	it.Loader = it
	it.IsZeroCenter = true
	it.Width, it.Height = 1024, 512
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
		dt = data.GetData(it.Sign, 0, 65536)
	} else if it.Sign != ViewSign {
		if dt = data.GetData(it.Sign, 0, 65536); dt != nil {
			dt.Sign = ViewSign
		}
	}
	if dt != nil {
		it.Load(dt)
		return true
	}
	return false
}

func (it *ItSignal) RunMouse(nb int, x, y float64, on bool) {
	val := it.X.ToVal(x)
	if nb == 1 && on {
		it.Croping = true
		it.CropFrom = val
		if it.CropTo < it.CropFrom {
			it.CropTo = it.X.Max
		}
	} else if nb == 2 && on {
		it.Croping = true
		it.CropTo = val
		if it.CropFrom > it.CropTo {
			it.CropFrom = it.X.Min
		}
	} else if nb == 3 && on {
		it.Croping = false
	}
	ViewSign++
}
