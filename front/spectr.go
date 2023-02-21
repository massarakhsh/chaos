package front

import (
	"math"

	"github.com/massarakhsh/chaos/data"

	"github.com/mjibson/go-dsp/fft"
)

const MAX_WAVE = 1024
const MAX_HIST = 16

type ItSpectr struct {
	ItPlot
	history [MAX_HIST][MAX_WAVE]float64
}

func BuildSpectr() *ItSpectr {
	graph := &ItSpectr{}
	graph.Name = "spectr"
	graph.Loader = graph
	graph.IsZeroCenter = false
	graph.Width, graph.Height = 512, 256
	return graph
}

func (it *ItSpectr) Probe() bool {
	if sign, step, vals := it.loadData(); vals == nil {
		return false
	} else {
		spectr := fft.FFTReal(vals)
		it.storeData(sign, step, spectr)
		return true
	}
}

func (it *ItSpectr) loadData() (int, float64, []float64) {
	if dt := data.GetData(it.Sign, 4096); dt == nil {
		return 0, 0, nil
	} else if count := dt.Length; count == 0 {
		return 0, 0, nil
	} else {
		sign := dt.Sign
		step := (dt.XMax - dt.XMin) / float64(dt.Length)
		vals := make([]float64, count)
		for n := 0; n < count; n++ {
			vals[n] = dt.Data[n]
		}
		return sign, step, vals
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
