package data

import (
	"fmt"
	"time"

	"github.com/tarm/serial"
)

var serInit bool
var SConf *serial.Config
var SPort *serial.Port

func serialInit() bool {
	for ns := 0; ns < 4; ns++ {
		dev := fmt.Sprintf("/dev/ttyUSB%d", ns)
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

func genPotSerial() []ItPot {
	var pots []ItPot
	if !serInit {
		serInit = serialInit()
	}
	if serInit {
		buf := make([]byte, 1)
		/*for nb := 0; nb < 100; nb++ {
			buf[0] = byte(nb)
			SPort.Write(buf)
		}*/
		for {
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
