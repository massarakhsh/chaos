package front

import (
	"math"
	"time"

	"github.com/andlabs/ui"
	"github.com/massarakhsh/chaos/data"
	"github.com/massarakhsh/chaos/pkg/zone"

	"github.com/mjibson/go-dsp/fft"
)

const MAX_FREQUENCY = 5000

type ItSpectr struct {
	zone.ItZone
	ItPlot

	Interval *ItInterval
	ViewSign int
}

func BuildSpectr(interval *ItInterval) *ItSpectr {
	it := &ItSpectr{Interval: interval}
	it.area = ui.NewArea(it)
	it.Mouse = it
	it.BindControl(it, it.area)
	it.Panel.Loader = it
	it.Panel.IsZeroCenter = false
	it.Panel.Width, it.Panel.Height = 512, 256
	it.BindRefresh(it)
	return it
}

func (it *ItSpectr) Refresh() {
	if time.Since(it.lastUpdate) >= time.Second*1 {
		if it.Probe() {
			it.area.QueueRedrawAll()
		}
		it.lastUpdate = time.Now()
	}
}

func (it *ItSpectr) Probe() bool {
	if it.Panel.Sign == it.Interval.Panel.Sign && it.ViewSign == it.Interval.ViewSign {
		return false
	} else if dt := it.Interval.Panel.GetData(); dt == nil {
		return false
	} else if length := len(dt.Data); length < 2 {
		return false
	} else {
		it.ViewSign = it.Interval.ViewSign
		sign := dt.Sign
		step := (dt.XMax - dt.XMin) / float64(length)
		vals := make([]float64, length)
		for n := 0; n < length; n++ {
			vals[n] = dt.Data[n]
		}
		spectr := fft.FFTReal(vals)
		it.storeData(sign, step, spectr)
		return true
	}
}

func (it *ItSpectr) storeData(sign int, step float64, info []complex128) {
	dia := len(info)
	zone := float64(dia) * step
	serial := &data.ItData{}
	serial.Sign = sign
	serial.XMin = 0
	serial.XMax = float64(MAX_FREQUENCY)
	serial.Data = make([]float64, MAX_FREQUENCY)
	left := 1.0
	ileft := 1
	for n := 1; n < MAX_FREQUENCY; n++ {
		summa := 0.0
		weight := 0.0
		frq := float64(n) * zone
		ifrq := int(math.Floor(frq))
		for ileft < ifrq && ileft < dia {
			doze := float64(ileft+1) - left
			summa += math.Abs(real(info[ileft])) * doze
			weight += doze
			ileft++
			left = float64(ileft)
		}
		if ileft < dia {
			doze := frq - left
			summa += math.Abs(real(info[ileft])) * doze
			weight += doze
			left = frq
		}
		if weight > 0 {
			serial.Data[n] = summa / weight
		}
	}
	it.Panel.Load(serial)
}
