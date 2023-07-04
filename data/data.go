package data

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
	//"github.com/gotk3/gotk3/gtk"
)

type ItData struct {
	Sign       string
	XMin, XMax float64
	Data       []float64
}

type ItPot struct {
	At   time.Time
	Data float64
}

const SOURCE_RESET = 0
const SOURCE_SERIAL = 1
const SOURCE_MODEL = 2
const SOURCE_FILE = 3
const SOURCE_ANALIZE = 4

var dataSource = SOURCE_RESET

var dataSpace = 1000000
var dataStart time.Time
var dataPots []ItPot
var dataFrom int
var dataTo int
var dataSize int
var dataLock sync.Mutex
var dataSign string

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

func GetSource() int {
	return dataSource
}

func SetSource(source int) {
	if source != SOURCE_ANALIZE {
		GenReset()
	}
	dataSource = source
}

func GetDataSize() int {
	return dataSize
}

func GenReset() {
	dataLock.Lock()
	genReset()
	dataLock.Unlock()
}

func genReset() {
	if dataSource == SOURCE_SERIAL {
		serialClose()
	} else if dataSource == SOURCE_MODEL {
		modelClose()
	}
	dataPots = make([]ItPot, dataSpace)
	dataFrom = 0
	dataTo = 0
	dataSize = 0
	dataSign = fmt.Sprint(rand.Intn(1000))
}

func genAppendPot(pots []ItPot) {
	dataLock.Lock()
	for _, pot := range pots {
		pushPot(pot)
	}
	dataSign = fmt.Sprint(rand.Intn(1000))
	dataLock.Unlock()
}

func pushPot(pot ItPot) {
	if dataSize == dataSpace {
		dataFrom = nextSource(dataFrom)
		dataSize--
	}
	dataPots[dataTo] = pot
	dataTo = nextSource(dataTo)
	dataSize++
}

func nextSource(pos int) int {
	if next := pos + 1; next < dataSpace {
		return next
	} else {
		return 0
	}
}

func SaveToFile() {
	if dataSize > 0 {
		filename := time.Now().Format("cha/2006-01-02-15-04-05.cha")
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			os.Mkdir("cha", 0777)
			file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
		}
		if err != nil {
			fmt.Println(err)
		} else {
			started := false
			var startAt time.Time
			for pos, size := dataFrom, dataSize; size > 0; pos, size = nextSource(pos), size-1 {
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

func GetData(sign string, first int, count int) *ItData {
	if sign == dataSign {
		return nil
	}
	dataLock.Lock()
	data := &ItData{}
	data.Sign = dataSign
	from := dataFrom
	length := dataSize
	if first > length {
		length = 0
	} else if first > 0 {
		length -= first
		from += first
		if from >= dataSpace {
			from -= dataSpace
		}
	}
	if length > count {
		from += length - count
		if from >= dataSpace {
			from -= dataSpace
		}
		length = count
	}
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

func ReadFromFile(filename string) {
	dataLock.Lock()
	genReset()
	if file, err := os.Open(filename); err == nil {
		now := time.Now()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var at, dt float64
			if no, err := fmt.Sscanf(scanner.Text(), "%f,%f", &at, &dt); no == 2 && err == nil {
				mcs := int64(at * 1000000)
				pot := ItPot{At: now.Add(time.Microsecond * time.Duration(mcs)), Data: dt}
				pushPot(pot)
			}
		}
		file.Close()
	}
	dataSign = fmt.Sprint(rand.Intn(1000))
	dataLock.Unlock()
}
