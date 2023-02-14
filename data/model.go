package data

import (
	"math"
	"time"
)

var modInit bool
var modNext time.Time
var modStep = time.Microsecond * 100

func genPotModel() []ItPot {
	var pots []ItPot
	now := time.Now()
	if !modInit {
		modNext = now
		modInit = true
	}
	for !modNext.After(now) {
		temp := modNext.Sub(dataStart).Seconds()
		value := 1.0*math.Sin(temp*math.Pi*41) + 1.0*math.Sin(temp*math.Pi*59)
		pot := ItPot{At: modNext, Data: value}
		pots = append(pots, pot)
		modNext = modNext.Add(modStep)
	}
	return pots
}
