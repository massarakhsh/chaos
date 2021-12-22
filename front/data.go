package front

import (
	"github/massarakhsh/chaos/data"

	"github.com/andlabs/ui"
)

type ItData struct {
	Sign          int
	SP            *ui.DrawStrokeParams
	Width, Height float64
	Count         int
	List          []ItPoint
	XMin, XMax    float64
	YMin, YMax    float64
	XZero, YZero  float64
}

type ItPoint struct {
	XVal, YVal float64
	XLoc, YLoc float64
}

func (it *ItData) Load(serial *data.ItSerial) {
	if serial == nil || serial.Length < 2 {
		it.Count = 0
		it.List = []ItPoint{}
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
		if it.YMin >= it.YMax {
			it.YMin -= 0.1
			it.YMax += 0.1
		}
	}
}

func (it *ItData) Calc() {
	it.XZero, it.YZero = 0, it.Height
	if length := it.Count; length >= 2 {
		x0 := it.XMin * it.Width / (it.XMin - it.XMax)
		if x0 < 0 {
			x0 = 0
		} else if x0 > it.Width {
			x0 = it.Width
		}
		it.XZero = x0
		y0 := it.YMax * it.Height / (it.YMax - it.YMin)
		if y0 < 0 {
			y0 = 0
		} else if y0 > it.Height {
			y0 = it.Height
		}
		it.YZero = y0
		for i := 0; i < length; i++ {
			point := &it.List[i]
			point.XLoc = (0*(it.XMax-point.XVal) + it.Width*(point.XVal-it.XMin)) / (it.XMax - it.XMin)
			point.YLoc = (0*(point.YVal-it.YMin) + it.Height*(it.YMax-point.YVal)) / (it.YMax - it.YMin)
		}
	}
}
