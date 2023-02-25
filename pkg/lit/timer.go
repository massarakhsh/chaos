package lit

import "time"

type ItTimer struct {
	lastAt time.Time
}

func (it *ItTimer) Probe(step time.Duration, view time.Duration) bool {
	yes := false
	if time.Since(it.lastAt) >= view {
		yes = true
		it.lastAt = time.Now()
	}
	time.Sleep(step)
	return yes
}
