package front

import "github/massarakhsh/chaos/data"

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
	if data.Data.Sign != it.Sign {
		it.Load(data.Data)
		return true
	} else {
		return false
	}
}
