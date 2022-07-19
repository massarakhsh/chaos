package front

import (
	"github/massarakhsh/chaos/data"
	"math"

	"github.com/mjibson/go-dsp/fft"
)

const MAX_WAVE = 256
const MAX_HIST = 16

type ItSpectr struct {
	ItSerial
	history [MAX_HIST][MAX_WAVE]float64
}

var MainSpectr ItSpectr

func BuildSpectr() *ItSpectr {
	MainSpectr.Loader = &MainSpectr
	MainGraphic.IsZeroMin = true
	return &MainSpectr
}

func (it *ItSpectr) Probe() bool {
	if MainGraphic.Sign != it.Sign {
		sign, data := it.loadData()
		spectr := fft.FFTReal(data)
		it.storeData(sign, spectr)
		return true
	} else {
		return false
	}
}

func (it *ItSpectr) loadData() (int, []float64) {
	sign := MainGraphic.Sign
	count := MainGraphic.Count
	info := make([]float64, count)
	for n := 0; n < count; n++ {
		info[n] = MainGraphic.List[n].YVal
	}
	return sign, info
}

func (it *ItSpectr) storeData(sign int, info []complex128) {
	for h := 0; h+1 < MAX_HIST; h++ {
		it.history[h] = it.history[h+1]
	}
	dia := len(info)
	serial := &data.ItData{}
	serial.Sign = sign
	serial.Length = MAX_WAVE
	serial.XMin = 0
	serial.XMax = float64(MAX_WAVE)
	serial.Data = make([]float64, MAX_WAVE)
	for n := 0; n < MAX_WAVE; n++ {
		ampl := 0.0
		if n > 1 {
			frq := float64(dia) / float64(n)
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
