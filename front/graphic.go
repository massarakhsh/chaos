package front

import "github/massarakhsh/chaos/data"

type ItGraphic struct {
	Data ItData
}

var graphic ItGraphic

func BuildGraphic() ItSource {
	return &graphic
}

func (it *ItGraphic) GetData() *ItData {
	return &it.Data
}
func (it *ItGraphic) LoadData() {
	it.Data.Load(data.Serial)
}
