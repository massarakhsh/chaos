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

	Signal *ItSignal
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
	if dt := it.getData(); dt != nil {
		it.Load(dt)
		return true
	} else {
		return false
	}
}

func (it *ItInterval) getData() *data.ItData {
	// if sign == dataSign {
	// 	return nil
	// }
	sigdata := it.Signal.Data
	if len(sigdata) < 16 {
		return nil
	}
	data := &data.ItData{}
	data.Sign = it.Signal.Sign
	first := 0
	length := len(sigdata)
	if length > 256 {
		first = length - 256
		length = 256
	}
	data.Length = length
	data.XMin = (it.Signal.X.Min*float64(length-first) + it.Signal.X.Max*float64(first)) / float64(length)
	data.XMax = it.Signal.X.Max
	data.Data = sigdata[first:]
	return data
}
