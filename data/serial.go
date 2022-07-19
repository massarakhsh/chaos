package data

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tarm/serial"
)

var SConf *serial.Config
var SPort *serial.Port

func InitSerial() {
	for ns := 0; ns < 4; ns++ {
		dev := fmt.Sprintf("/dev/ttyUSB%d", ns)
		conf := &serial.Config{Name: dev, Baud: 115200, ReadTimeout: time.Millisecond * 10}
		if port, err := serial.OpenPort(conf); err == nil {
			SConf = conf
			SPort = port
			break
		}
	}
}

func LoadSerial() {
	if SPort == nil {
		InitSerial()
		if SPort == nil {
			return
		}
	}
	buf := make([]byte, 1)
	if sbuf, err := SPort.Read(buf); err != nil || sbuf == 0 {
		return
	} else {
		for i := 0; i+1 < Data.Length; i++ {
			Data.Data[i] = Data.Data[i+1]
		}
		Data.Data[Data.Length-1] = float64(int(buf[0]) - 128)
		Data.Sign = rand.Int()
	}
}
