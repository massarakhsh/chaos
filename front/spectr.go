package front

import (
	"github/massarakhsh/chaos/data"
	"math"

	"github.com/mjibson/go-dsp/fft"
)

type ItSpectr struct {
	ItSerial
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
	count := 512
	dia := len(info)
	serial := &data.ItData{}
	serial.Sign = sign
	serial.Length = count
	serial.XMin = 0
	serial.XMax = float64(count)
	serial.Data = make([]float64, count)
	for n := 0; n < count; n++ {
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
		serial.Data[n] = ampl
	}
	it.Load(serial)
}
