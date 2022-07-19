package front

import (
	"github.com/massarakhsh/chaos/data"
)

type ItGraphic struct {
	ItSerial
}

var MainGraphic ItGraphic

func BuildGraphic() *ItGraphic {
	MainGraphic.Loader = &MainGraphic
	MainGraphic.IsZeroCenter = true
	return &MainGraphic
}

func (it *ItGraphic) Probe() bool {
	if dt := data.GetData(it.Sign); dt != nil {
		it.Load(dt)
		return true
	} else {
		return false
	}
}
