package data

import "math"

type ItSerial struct {
	Length int
	Data   []float64
}

var Serial ItSerial

func Generate() {
	Serial = ItSerial{}
	Serial.Length = 1024
	Serial.Data = make([]float64, Serial.Length)
	for n := 0; n < Serial.Length; n++ {
		Serial.Data[n] = float64(n*n) / 1024 / 1024 * math.Sin(float64(n)/100)
	}
}
