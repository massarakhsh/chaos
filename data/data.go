package data

import (
	"math"
	"time"
)

type ItSerial struct {
	Sign       int
	Length     int
	XMin, XMax float64
	Data       []float64
}

var Sign = 0
var Serial *ItSerial

func Generate() {
	generate()
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			generate()
		}
	}()
}

func generate() {
	Sign++
	length := 1024
	if Serial != nil {
		length = Serial.Length + 1
	}
	serial := &ItSerial{}
	serial.Sign = Sign
	serial.Length = length
	serial.XMin = float64(-(length / 100))
	serial.XMax = serial.XMin + float64(length)
	serial.Data = make([]float64, serial.Length)
	for n := 0; n < serial.Length; n++ {
		serial.Data[n] = float64(n*n) / float64(length*length) * math.Sin(float64(n)/100)
	}
	Serial = serial
}
