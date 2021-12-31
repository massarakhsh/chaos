package data

import (
	"math"
	"math/rand"
	"time"
)

type ItData struct {
	Sign       int
	Length     int
	XMin, XMax float64
	Data       []float64
}

var Data *ItData

func Generate() {
	genInit()
	InitSerial()
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			LoadSerial()
		}
	}()
}

func genInit() {
	length := 1024
	if Data != nil {
		length = Data.Length + 1
	}
	serial := &ItData{}
	serial.Sign = rand.Int()
	serial.Length = length
	serial.XMin = float64(-(length / 100))
	serial.XMax = serial.XMin + float64(length)
	serial.Data = make([]float64, serial.Length)
	for n := 0; n < serial.Length; n++ {
		serial.Data[n] = 0.8*math.Sin(float64(n)*math.Pi/100) + 1.2*math.Cos(float64(n)*math.Pi/70)
	}
	Data = serial
}
