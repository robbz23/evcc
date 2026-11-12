package core

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/vehicle"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHoldWindow(t *testing.T) {
	config.Reset()
	t.Cleanup(config.Reset)

	ctrl := gomock.NewController(t)
	v := api.NewMockVehicle(ctrl)
	v.EXPECT().OnIdentified().Return(api.ActionConfig{}).AnyTimes()

	const name = "niro"
	require.NoError(t, config.Vehicles().Add(
		config.NewStaticDevice(config.Named{Name: name}, api.Vehicle(v)),
	))

	// 2026-11-02 is a Monday
	clk := clock.NewMock()
	clk.Set(time.Date(2026, 11, 2, 5, 0, 0, 0, time.UTC))

	lp := NewLoadpoint(util.NewLogger("foo"), nil)
	lp.clock = clk
	lp.vehicle = v
	lp.minCurrent, lp.maxCurrent = 6, 16

	// no windows
	_, ok := lp.activeHoldWindow()
	assert.False(t, ok)

	require.NoError(t, vehicle.Settings(lp.log, v).SetHoldWindows([]api.HoldWindow{{
		Weekdays: []int{1}, Start: "06:00", End: "07:30", Tz: "UTC",
		MinCurrent: 8, MaxCurrent: 10, Active: true,
	}}))

	// before the window
	_, ok = lp.activeHoldWindow()
	assert.False(t, ok)

	// inside the window
	clk.Set(time.Date(2026, 11, 2, 6, 30, 0, 0, time.UTC))
	w, ok := lp.activeHoldWindow()
	require.True(t, ok)

	minCurrent, maxCurrent := lp.holdWindowCurrents(w)
	assert.Equal(t, 8.0, minCurrent)
	assert.Equal(t, 10.0, maxCurrent)

	// no surplus: never disabled, held at window minimum
	assert.Equal(t, 8.0, lp.holdWindowCurrent(w, 0))
	// moderate surplus passes through
	assert.Equal(t, 9.0, lp.holdWindowCurrent(w, 9))
	// large surplus is capped at window maximum
	assert.Equal(t, 10.0, lp.holdWindowCurrent(w, 16))

	// after the window
	clk.Set(time.Date(2026, 11, 2, 7, 30, 0, 0, time.UTC))
	_, ok = lp.activeHoldWindow()
	assert.False(t, ok)
}

func TestHoldWindowDefaultCurrents(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := api.NewMockVehicle(ctrl)
	v.EXPECT().OnIdentified().Return(api.ActionConfig{}).AnyTimes()

	lp := NewLoadpoint(util.NewLogger("foo"), nil)
	lp.vehicle = v
	lp.minCurrent, lp.maxCurrent = 6, 16

	// zero current limits fall back to the loadpoint's effective range
	w := api.HoldWindow{}
	minCurrent, maxCurrent := lp.holdWindowCurrents(w)
	assert.Equal(t, 6.0, minCurrent)
	assert.Equal(t, 16.0, maxCurrent)

	// window limits can't go below the loadpoint minimum or above its maximum
	w = api.HoldWindow{MinCurrent: 2, MaxCurrent: 32}
	minCurrent, maxCurrent = lp.holdWindowCurrents(w)
	assert.Equal(t, 6.0, minCurrent)
	assert.Equal(t, 16.0, maxCurrent)
}
