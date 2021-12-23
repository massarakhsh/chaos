package front

import "github/massarakhsh/chaos/data"

type ItGraphic struct {
	ItPlot
}

var graphic ItGraphic

func BuildGraphic() *ItGraphic {
	graphic := &ItGraphic{}
	graphic.Source = graphic
	return graphic
}

func (it *ItGraphic) LoadData() {
	it.Load(data.Serial)
}
