package data

import (
	"sync"
	"time"
)

type ItData struct {
	Sign       int
	Length     int
	XMin, XMax float64
	Data       []float64
}

type ItPot struct {
	At   time.Time
	Data float64
}

var sourceDir = 0
var sourceLength = 16394
var sourceStart time.Time
var sourceData []ItPot
var sourcePos int
var sourceSign int
var sourceLock sync.Mutex

func StartData() {
	sourceStart = time.Now()
	GenReset()
	go func() {
		for {
			var pots []ItPot
			if sourceDir == 0 {
				pots = genPotModel()
			}
			if pots != nil {
				genAppendPot(pots)
			}
			time.Sleep(time.Microsecond * 1000)
		}
	}()
}

func GenReset() {
	sourceLock.Lock()
	sourceData = make([]ItPot, sourceLength)
	for old := 0; old < sourceLength; old++ {
		at := sourceStart.Add(-time.Duration(old) * time.Microsecond * 100)
		sourceData[sourceLength-1-old].At = at
		sourceData[sourceLength-1-old].Data = 0
	}
	sourcePos = 0
	sourceSign++
	sourceLock.Unlock()
}

func genAppendPot(pots []ItPot) {
	sourceLock.Lock()
	for nr := 0; nr < len(pots); nr++ {
		sourceData[sourcePos] = pots[nr]
		sourcePos++
		if sourcePos >= sourceLength {
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
	data.Length = sourceLength
	data.Data = make([]float64, sourceLength)
	for to := 0; to < sourceLength; to++ {
		from := sourcePos + to
		if from >= sourceLength {
			from -= sourceLength
		}
		data.Data[to] = sourceData[from].Data
		if to == 0 {
			data.XMin = sourceData[from].At.Sub(sourceStart).Seconds()
		}
		if to == sourceLength-1 {
			data.XMax = sourceData[from].At.Sub(sourceStart).Seconds()
		}
	}
	sourceLock.Unlock()
	return data
}
