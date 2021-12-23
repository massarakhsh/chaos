package front

import (
	"github/massarakhsh/chaos/data"

	"github.com/mjibson/go-dsp/fft"
)

type ItSpectr struct {
	ItSerial
}

var MainSpectr ItSpectr

func BuildSpectr() *ItSpectr {
	MainSpectr.Loader = &MainSpectr
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
	count := 40
	serial := &data.ItData{}
	serial.Sign = sign
	serial.Length = count
	serial.XMin = 0
	serial.XMax = float64(count)
	serial.Data = make([]float64, count)
	for n := 0; n < count; n++ {
		serial.Data[n] = real(info[n])
	}
	it.Load(serial)
}
