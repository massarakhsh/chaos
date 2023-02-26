package data

import (
	"fmt"
	"runtime"
	"time"

	"github.com/tarm/serial"
)

var serInit bool
var SConf *serial.Config
var SPort *serial.Port

func serialInit() bool {
	serialClose()
	for ns := 7; ns >= 0; ns-- {
		var dev string
		if os := runtime.GOOS; os == "linux" {
			dev = fmt.Sprintf("/dev/ttyUSB%d", ns)
		} else if os == "windows" {
			dev = fmt.Sprintf("COM%d", ns)
		}
		conf := &serial.Config{Name: dev, Baud: 115200, ReadTimeout: time.Millisecond * 10}
		if port, err := serial.OpenPort(conf); err == nil {
			SConf = conf
			SPort = port
			fmt.Printf("Found serial port %s\n", SConf.Name)
			return true
		}
	}
	return false
}

func serialClose() {
	if SPort != nil {
		SPort.Close()
		SPort = nil
		fmt.Printf("Closed serial port %s\n", SConf.Name)
	}
	SConf = nil
}

func genPotSerial() []ItPot {
	var pots []ItPot
	if !serInit {
		serInit = serialInit()
	}
	if serInit {
		buf := make([]byte, 1)
		for nn := 0; nn < 256; nn++ {
			if sbuf, err := SPort.Read(buf); err != nil || sbuf == 0 {
				break
			} else {
				value := float64(buf[0])
				pot := ItPot{At: time.Now(), Data: value}
				pots = append(pots, pot)
			}
		}
	}
	return pots
}
