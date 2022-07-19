package data

import (
	"math"
	"sync"
	"time"
)

type ItData struct {
	Sign       int
	Length     int
	XMin, XMax float64
	Data       []float64
}

type ItReal struct {
	At   time.Time
	Data float64
}

const sourceMax = 4096

var sourceStart time.Time
var sourceData []ItReal
var sourcePos int
var sourceSign int
var sourceLock sync.Mutex

func Generate() {
	genPreset()
	go func() {
		for {
			genAppendMath()
			time.Sleep(time.Microsecond * 1000)
		}
	}()
}

func genPreset() {
	sourceLock.Lock()
	sourceStart = time.Now()
	sourceData = make([]ItReal, sourceMax)
	for old := 0; old < sourceMax; old++ {
		at := sourceStart.Add(-time.Duration(old) * time.Microsecond)
		sourceData[sourceMax-1-old].At = at
		sourceData[sourceMax-1-old].Data = 0
	}
	sourcePos = 0
	sourceSign++
	sourceLock.Unlock()
}

func genAppendMath() {
	var real ItReal
	real.At = time.Now()
	temp := real.At.Sub(sourceStart).Seconds()
	value := 1.0*math.Sin(temp*math.Pi*20) + 0.5*math.Sin(temp*math.Pi*41)
	real.Data = value
	genAppendReal([]ItReal{real})
}

func genAppendReal(reals []ItReal) {
	sourceLock.Lock()
	for nr := 0; nr < len(reals); nr++ {
		sourceData[sourcePos] = reals[nr]
		sourcePos++
		if sourcePos >= sourceMax {
			sourcePos = 0
		}
	}
	sourceSign++
	sourceLock.Unlock()
}

func GetData(sign int) *ItData {
	if sign == sourceSign {
		return nil
	}
	sourceLock.Lock()
	data := &ItData{}
	data.Sign = sourceSign
	data.Length = sourceMax
	data.Data = make([]float64, sourceMax)
	for old := 0; old < sourceMax; old++ {
		pos := sourcePos + old
		if pos >= sourceMax {
			pos -= sourceMax
		}
		data.Data[old] = sourceData[pos].Data
		if old == 0 {
			data.XMin = sourceData[pos].At.Sub(sourceStart).Seconds()
		}
		if old == sourceMax-1 {
			data.XMax = sourceData[pos].At.Sub(sourceStart).Seconds()
		}
	}
	sourceLock.Unlock()
	return data
}
