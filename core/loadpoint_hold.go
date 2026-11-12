package core

import (
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/vehicle"
)

// activeHoldWindow returns the vehicle's currently active hold window, if any
func (lp *Loadpoint) activeHoldWindow() (api.HoldWindow, bool) {
	v := lp.GetVehicle()
	if v == nil {
		return api.HoldWindow{}, false
	}

	now := lp.clock.Now()
	for _, w := range vehicle.Settings(lp.log, v).GetHoldWindows() {
		if w.ActiveAt(now) {
			return w, true
		}
	}

	return api.HoldWindow{}, false
}

// holdWindowCurrents returns the current range of a hold window, bounded by the effective limits
func (lp *Loadpoint) holdWindowCurrents(w api.HoldWindow) (minCurrent, maxCurrent float64) {
	effMin, effMax := lp.effectiveMinCurrent(), lp.effectiveMaxCurrent()

	minCurrent = max(w.MinCurrent, effMin)
	maxCurrent = effMax
	if w.MaxCurrent > 0 {
		maxCurrent = min(w.MaxCurrent, effMax)
	}

	return minCurrent, max(minCurrent, maxCurrent)
}

// holdWindowCurrent bounds the target current to the active hold window.
// The charger is never disabled and never exceeds the window's maximum.
func (lp *Loadpoint) holdWindowCurrent(w api.HoldWindow, target float64) float64 {
	minCurrent, maxCurrent := lp.holdWindowCurrents(w)
	return min(max(target, minCurrent), maxCurrent)
}
