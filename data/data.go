package data

import (
	"math"
	"time"
)

type ItData struct {
	Sign       int
	Length     int
	XMin, XMax float64
	Data       []float64
}

var SignData = 0
var Data *ItData

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
	SignData++
	length := 1024
	if Data != nil {
		length = Data.Length + 1
	}
	serial := &ItData{}
	serial.Sign = SignData
	serial.Length = length
	serial.XMin = float64(-(length / 100))
	serial.XMax = serial.XMin + float64(length)
	serial.Data = make([]float64, serial.Length)
	for n := 0; n < serial.Length; n++ {
		serial.Data[n] = 0.8*math.Sin(float64(n)/100) + 1.2*math.Cos(float64(n)/70)
	}
	Data = serial
}
