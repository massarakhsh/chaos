package data

import (
	"math"
	"time"
)

var modInit bool
var modNext time.Time
var modStep = time.Microsecond * 100
var freq1 = 300.0
var freq2 = 511.0

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
		value := 1.0*math.Sin(temp*freq1*math.Pi*2) + 1.0*math.Sin(temp*freq2*math.Pi*2)
		pot := ItPot{At: modNext, Data: value}
		pots = append(pots, pot)
		modNext = modNext.Add(modStep)
	}
	return pots
}
