package front

import (
	"math"
	"time"

	"github.com/andlabs/ui"
	"github.com/massarakhsh/chaos/data"
	"github.com/massarakhsh/chaos/pkg/zone"

	"github.com/mjibson/go-dsp/fft"
)

const MAX_WAVE = 1024
const MAX_HIST = 16

type ItSpectr struct {
	zone.ItZone
	ItPlot
	viewSign int
	history  [MAX_HIST][MAX_WAVE]float64
}

func BuildSpectr() *ItSpectr {
	it := &ItSpectr{}
	it.area = ui.NewArea(it)
	it.BindControl(it, it.area)
	it.Loader = it
	it.IsZeroCenter = false
	it.Width, it.Height = 512, 256
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
	if dt := data.GetData(it.Sign, 0, 65536); dt == nil || dt.Length < 2 {
		return false
	} else {
		sign := dt.Sign
		step := (dt.XMax - dt.XMin) / float64(dt.Length)
		vals := make([]float64, dt.Length)
		for n := 0; n < dt.Length; n++ {
			vals[n] = dt.Data[n]
		}
		spectr := fft.FFTReal(vals)
		it.storeData(sign, step, spectr)
		return true
	}
}

func (it *ItSpectr) storeData(sign int, step float64, info []complex128) {
	for h := 0; h+1 < MAX_HIST; h++ {
		it.history[h] = it.history[h+1]
	}
	dia := len(info)
	zone := float64(dia) * step
	serial := &data.ItData{}
	serial.Sign = sign
	serial.Length = MAX_WAVE
	serial.XMin = 0
	serial.XMax = float64(MAX_WAVE)
	serial.Data = make([]float64, MAX_WAVE)
	for n := 0; n < MAX_WAVE; n++ {
		ampl := 0.0
		if n > 1 {
			frq := float64(n) * zone
			rfrq := math.Floor(frq)
			ifrq := int(rfrq)
			if ifrq > 0 && ifrq+1 < dia {
				left := math.Abs(real(info[ifrq]))
				right := math.Abs(real(info[ifrq+1]))
				ampl = left*(rfrq+1-frq) + right*(frq-rfrq)
			}
		}
		it.history[MAX_HIST-1][n] = ampl
	}
	for n := 0; n < MAX_WAVE; n++ {
		ampl := 0.0
		for h := 0; h < MAX_HIST; h++ {
			ampl += it.history[h][n]
		}
		serial.Data[n] = ampl / MAX_HIST
	}
	it.Load(serial)
}
