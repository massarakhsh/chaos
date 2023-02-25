package front

import (
	"time"

	"github.com/andlabs/ui"
)

func (it *ItFront) setInfoMode() {
	it.purgeInfo()
	it.buildInfoBox()
}

func (it *ItFront) purgeInfo() {
	if it.RowBox != nil {
		for na := len(it.ListAreas) - 1; na >= 0; na-- {
			it.ListAreas[na].Disable()
			time.Sleep(time.Millisecond * 20)
			it.RowBox.Delete(na)
		}
		it.ListAreas = []*ui.Area{}
		it.InfoBox.Delete(1)
		it.RowBox = nil
	}
	if it.InfoArea != nil {
		it.InfoArea.Disable()
		time.Sleep(time.Millisecond * 20)
		it.InfoBox.Delete(0)
		it.InfoArea = nil
	}
}

func (it *ItFront) buildInfoBox() {
	it.InfoArea = it.buildPlotArea(BuildGraphic(it))
	it.InfoBox.Append(it.InfoArea, true)
	it.RowBox = ui.NewHorizontalBox()
	it.RowBox.SetPadded(true)
	if elm := it.buildPlotArea(BuildSignal(it)); elm != nil {
		it.RowBox.Append(elm, true)
		it.ListAreas = append(it.ListAreas, elm)
	}
	if elm := it.buildPlotArea(BuildSpectr(it)); elm != nil {
		it.RowBox.Append(elm, true)
		it.ListAreas = append(it.ListAreas, elm)
	}
	it.InfoBox.Append(it.RowBox, true)
}

func (it *ItFront) buildPlotArea(plot InfPlot) *ui.Area {
	area := ui.NewArea(plot)
	go func() {
		nexttime := time.Now()
		for area.Enabled() {
			time.Sleep(time.Millisecond * 10)
			if time.Now().After(nexttime) {
				if DuraUpdate > 0 {
					nexttime = time.Now().Add(time.Millisecond * time.Duration(DuraUpdate))
				}
				if DuraUpdate > 0 || (NeedUpdate&0x1) != 0 {
					NeedUpdate &= 0xfe
					if plot.Probe() {
						area.QueueRedrawAll()
					}
				}
			}
		}
	}()
	return area
}
