package front

import (
	"fmt"
	"math"

	"github.com/andlabs/ui"
	"github.com/massarakhsh/chaos/data"
)

type IfLoad interface {
	Probe() bool
}

type ItAxis struct {
	Min, Max    float64
	First, Step float64

	LocDep  float64
	LocSize float64
	LocZero float64
	Acc     bool

	Format string
}

type ItPanel struct {
	Name          string
	Sign          int
	Width, Height float64
	IsZeroCenter  bool
	IsZeroMin     bool

	Data []float64
	X    ItAxis
	Y    ItAxis

	Croping  bool
	CropFrom float64
	CropTo   float64

	Loader IfLoad
}

func (it *ItAxis) Calibrate(dep, size float64, acc bool) {
	it.LocDep = dep
	it.LocSize = size
	it.Acc = acc
	it.LocZero = it.ToLoc(0)
	if it.Max <= it.Min {
		it.First = it.Min
		it.Step = 1.0
		it.Format = "%f"
	} else {
		step := 1.0
		dg := 0
		for step*10 > it.Max-it.Min {
			step /= 10
			dg++
		}
		for step*100 < it.Max-it.Min {
			step *= 10
			if dg > 0 {
				dg--
			}
		}
		for step*20 < it.Max-it.Min {
			step *= 2
		}
		first := math.Floor(it.Min/step)*step - step
		for first < it.Min {
			first += step
		}
		it.First = first
		it.Step = step
		it.Format = fmt.Sprintf("%%.%df", dg)
	}
}

func (it *ItAxis) ToLoc(val float64) float64 {
	if it.Max <= it.Min {
		return 0
	}
	loc := (0*(it.Max-val) + float64(it.LocSize)*(val-it.Min)) / (it.Max - it.Min)
	if loc < 0 {
		loc = 0
	} else if loc > it.LocSize-1 {
		loc = it.LocSize - 1
	}
	if !it.Acc {
		loc = it.LocSize - 1 - loc
	}
	return it.LocDep + loc
}

func (it *ItAxis) ToVal(loc float64) float64 {
	if it.LocSize <= 0 {
		return (it.Min + it.Max) / 2
	}
	val := (it.Max*loc + it.Min*(it.LocSize-loc)) / it.LocSize
	return val
}

func (it *ItPanel) Load(serial *data.ItData) {
	if serial == nil {
		it.Data = []float64{}
		it.X = ItAxis{}
		it.Y = ItAxis{}
	} else if length := serial.Length; length < 2 {
		it.Data = []float64{}
		it.X = ItAxis{}
		it.Y = ItAxis{}
	} else {
		it.Sign = serial.Sign
		it.Data = serial.Data
		it.X = ItAxis{Min: serial.XMin, Max: serial.XMax}
		if it.X.Min >= it.X.Max {
			mid := (it.X.Min + it.X.Max) / 2
			it.X.Min = mid - 0.1
			it.X.Max = mid + 0.1
		}
		it.Y = ItAxis{}
		for n, val := range it.Data {
			if n == 0 || val < it.Y.Min {
				it.Y.Min = val
			}
			if n == 0 || val > it.Y.Max {
				it.Y.Max = val
			}
		}
		if it.Y.Min >= it.Y.Max {
			it.Y.Min -= 0.1
			it.Y.Max += 0.1
		}
		if it.IsZeroCenter {
			if it.Y.Min >= 0 {
				it.Y.Min = -it.Y.Max
			} else if it.Y.Max <= 0 {
				it.Y.Max = -it.Y.Min
			} else if it.Y.Min > -it.Y.Max {
				it.Y.Min = -it.Y.Max
			} else if it.Y.Max < -it.Y.Min {
				it.Y.Max = -it.Y.Min
			}
		}
		if it.IsZeroMin {
			if it.Y.Min > 0 {
				it.Y.Min = 0
			}
		}
	}
}

func (it *ItPanel) resize(p *ui.AreaDrawParams) {
	it.Width = p.AreaWidth
	it.Height = p.AreaHeight
	m := ui.DrawNewMatrix()
	//m.Translate(bd, bd)
	p.Context.Transform(m)
}
