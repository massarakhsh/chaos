package data

import (
	"math"
	"time"
)

var modInit bool
var modNext time.Time
var modStep = time.Microsecond * 100
var freq1 = 180.0
var amp1 = 0.5
var freq2 = 373.0
var amp2 = 1.0
var freq3 = 505.0
var amp3 = 1.5

func modelClose() {
	modInit = false
}

func genPotModel() []ItPot {
	var pots []ItPot
	now := time.Now()
	if !modInit {
		modNext = now
		modInit = true
	}
	for !modNext.After(now) {
		temp := modNext.Sub(dataStart).Seconds()
		//value := 1.0*math.Sin(temp*math.Pi*41) + 1.0*math.Sin(temp*math.Pi*59)
		value := amp1*math.Sin(temp*freq1*math.Pi*2) + amp2*math.Sin(temp*freq2*math.Pi*2) + amp3*math.Sin(temp*freq3*math.Pi*2)
		pot := ItPot{At: modNext, Data: value}
		pots = append(pots, pot)
		modNext = modNext.Add(modStep)
	}
	return pots
}
