package front

import (
	"github.com/andlabs/ui"
	"github.com/massarakhsh/chaos/data"
)

type ItLoad interface {
	Probe() bool
}

type ItPanel struct {
	Name          string
	Sign          int
	Width, Height float64
	Count         int
	List          []ItPoint
	XMin, XMax    float64
	YMin, YMax    float64
	XZero, YZero  float64
	IsZeroCenter  bool
	IsZeroMin     bool
	XFirst, XStep float64
	Xfmt          string
	YFirst, YStep float64
	Yfmt          string

	Loader ItLoad
}

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
		it.XMin = serial.XMin
		it.XMax = serial.XMax
		if it.XMin >= it.XMax {
			it.XMin -= 0.1
			it.XMax += 0.1
		}
		for i := 0; i < length; i++ {
			point := &it.List[i]
			point.XVal = (serial.XMin*float64(length-1-i) + serial.XMax*float64(i)) / float64(length-1)
			val := serial.Data[i]
			point.YVal = val
			if i == 0 || val < it.YMin {
				it.YMin = val
			}
			if i == 0 || val > it.YMax {
				it.YMax = val
			}
		}
		if it.IsZeroCenter {
			if it.YMin >= 0 {
				it.YMin = -it.YMax
			} else if it.YMax <= 0 {
				it.YMax = -it.YMin
			} else if it.YMin > -it.YMax {
				it.YMin = -it.YMax
			} else if it.YMax < -it.YMin {
				it.YMax = -it.YMin
			}
		}
		if it.IsZeroMin {
			if it.YMin > 0 {
				it.YMin = 0
			}
		}
		if it.YMin >= it.YMax {
			it.YMin -= 0.1
			it.YMax += 0.1
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
