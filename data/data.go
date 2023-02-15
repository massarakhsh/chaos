package data

import (
	"fmt"
	"os"
	"sync"
	"time"
	//"github.com/gotk3/gotk3/gtk"
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

const SOURCE_NO = 0
const SOURCE_SERIAL = 1
const SOURCE_MODEL = 2

var dataSource = SOURCE_NO

var dataSpace = 1024 * 1024
var dataStart time.Time
var dataPots []ItPot
var dataFrom int
var dataTo int
var dataSign int
var dataLock sync.Mutex

func StartData() {
	dataStart = time.Now()
	GenReset()
	go func() {
		for {
			var pots []ItPot
			if dataSource == SOURCE_SERIAL {
				pots = genPotSerial()
			} else if dataSource == SOURCE_MODEL {
				pots = genPotModel()
			}
			if pots != nil {
				genAppendPot(pots)
			}
			time.Sleep(time.Microsecond * 100)
		}
	}()
}

func SetSource(source int) {
	dataSource = source
}

func GenReset() {
	dataLock.Lock()
	dataPots = make([]ItPot, dataSpace)
	dataFrom = 0
	dataTo = 0
	dataSign++
	dataLock.Unlock()
}

func genAppendPot(pots []ItPot) {
	dataLock.Lock()
	for nr := 0; nr < len(pots); nr++ {
		dataPots[dataTo] = pots[nr]
		dataTo = nextSource(dataTo)
		if dataTo == dataFrom {
			dataFrom = nextSource(dataFrom)
		}
	}
	dataSign++
	dataLock.Unlock()
}

func nextSource(pos int) int {
	if next := pos + 1; next < dataSpace {
		return next
	} else {
		return 0
	}
}

func SaveToFile() {
	if dataTo != dataFrom {
		filename := time.Now().Format("2006-01-02 15:04:05") + ".cha"
		file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
		if file != nil {
			started := false
			var startAt time.Time
			to := dataTo
			for pos := dataFrom; pos != to; pos = nextSource(pos) {
				pot := dataPots[pos]
				if !started {
					startAt = pot.At
					started = true
				}
				file.WriteString(fmt.Sprintf("%.6f,%.6f\n", pot.At.Sub(startAt).Seconds(), pot.Data))
			}
			file.Close()
		}
	}
}

func GetData(sign int, max int) *ItData {
	if sign == dataSign {
		return nil
	}
	dataLock.Lock()
	data := &ItData{}
	data.Sign = dataSign
	from := dataFrom
	length := dataTo - dataFrom
	if length < 0 {
		length += dataSpace
	}
	if length > max {
		from += length - max
		if from >= dataSpace {
			from -= dataSpace
		}
		length = max
	}
	data.Length = length
	data.Data = make([]float64, length)
	data.XMin = dataPots[from].At.Sub(dataStart).Seconds()
	for n := 0; n < length; n++ {
		data.XMax = dataPots[from].At.Sub(dataStart).Seconds()
		data.Data[n] = dataPots[from].Data
		from = nextSource(from)
	}
	dataLock.Unlock()
	return data
}
