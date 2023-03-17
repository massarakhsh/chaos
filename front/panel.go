package front

import (
	"fmt"
	"math"

	"github.com/andlabs/ui"
	"github.com/massarakhsh/chaos/data"
)

type ItFLoad interface {
	Probe() bool
}

type ItAxis struct {
	Min, Max float64
	LocDep   float64
	LocSize  float64
	LocZero  float64
	Acc      bool

	Format      string
	First, Step float64
}

type ItPanel struct {
	Front         *ItFront
	Name          string
	Sign          int
	Width, Height float64
	Count         int
	List          []ItPoint
	IsZeroCenter  bool
	IsZeroMin     bool

	X ItAxis
	Y ItAxis

	CropFrom float64
	CropTo   float64

	Loader ItFLoad
}

func (it *ItAxis) Calibrate(dep float64, size float64, acc bool) {
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

func (it *ItAxis) ToLoc(v float64) float64 {
	if it.Max <= it.Min {
		return 0
	}
	loc := (0*(it.Max-v) + it.LocSize*(v-it.Min)) / (it.Max - it.Min)
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

// func (it *ItPlot) locToX(x float64) float64 {
// 	return (it.XMin*(it.Width-x) + it.XMax*(x-0)) / it.Width
// }
// func (it *ItPlot) locToY(y float64) float64 {
// 	return (it.YMin*(y-0) + it.YMax*(it.Height-y)) / it.Height
// }

func (it *ItPanel) Load(serial *data.ItData) {
	if serial == nil {
		it.Count = 0
		it.List = []ItPoint{}
	} else if serial.Length < 2 {
		it.Count = 0
		it.List = []ItPoint{}
		it.Sign = serial.Sign
	} else {
		length := serial.Length
		it.Sign = serial.Sign
		it.Count = length
		it.List = make([]ItPoint, length)
		it.X = ItAxis{Min: serial.XMin, Max: serial.XMax}
		it.Y = ItAxis{}
		if it.X.Min >= it.X.Max {
			mid := (it.X.Min + it.X.Max) / 2
			it.X.Min = mid - 0.1
			it.X.Max = mid + 0.1
		}
		for i := 0; i < length; i++ {
			point := &it.List[i]
			point.XVal = (serial.XMin*float64(length-1-i) + serial.XMax*float64(i)) / float64(length-1)
			val := serial.Data[i]
			point.YVal = val
			if i == 0 || val < it.Y.Min {
				it.Y.Min = val
			}
			if i == 0 || val > it.Y.Max {
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
