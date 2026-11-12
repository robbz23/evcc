package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHoldWindowActiveAt(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)

	// 2026-11-02 is a Monday
	at := func(day, hour, min int) time.Time {
		return time.Date(2026, 11, day, hour, min, 0, 0, loc)
	}

	weekdays := []int{1, 2, 3, 4, 5} // Mon-Fri

	tc := []struct {
		name string
		w    HoldWindow
		ts   time.Time
		want bool
	}{
		{"inside", HoldWindow{Weekdays: weekdays, Start: "06:00", End: "07:30", Tz: "Europe/Berlin", Active: true}, at(2, 6, 45), true},
		{"start inclusive", HoldWindow{Weekdays: weekdays, Start: "06:00", End: "07:30", Tz: "Europe/Berlin", Active: true}, at(2, 6, 0), true},
		{"end exclusive", HoldWindow{Weekdays: weekdays, Start: "06:00", End: "07:30", Tz: "Europe/Berlin", Active: true}, at(2, 7, 30), false},
		{"before", HoldWindow{Weekdays: weekdays, Start: "06:00", End: "07:30", Tz: "Europe/Berlin", Active: true}, at(2, 5, 59), false},
		{"wrong weekday", HoldWindow{Weekdays: weekdays, Start: "06:00", End: "07:30", Tz: "Europe/Berlin", Active: true}, at(7, 6, 45), false}, // Saturday
		{"inactive", HoldWindow{Weekdays: weekdays, Start: "06:00", End: "07:30", Tz: "Europe/Berlin", Active: false}, at(2, 6, 45), false},
		{"no weekdays", HoldWindow{Start: "06:00", End: "07:30", Tz: "Europe/Berlin", Active: true}, at(2, 6, 45), false},
		{"timezone converted", HoldWindow{Weekdays: weekdays, Start: "06:00", End: "07:30", Tz: "Europe/Berlin", Active: true}, at(2, 6, 45).UTC(), true},
		{"midnight evening part", HoldWindow{Weekdays: []int{5}, Start: "22:00", End: "02:00", Tz: "Europe/Berlin", Active: true}, at(6, 23, 0), true},      // Fri 23:00
		{"midnight morning part", HoldWindow{Weekdays: []int{5}, Start: "22:00", End: "02:00", Tz: "Europe/Berlin", Active: true}, at(7, 1, 0), true},       // Sat 01:00
		{"midnight morning other day", HoldWindow{Weekdays: []int{5}, Start: "22:00", End: "02:00", Tz: "Europe/Berlin", Active: true}, at(6, 1, 0), false}, // Fri 01:00
		{"midnight gap", HoldWindow{Weekdays: []int{5}, Start: "22:00", End: "02:00", Tz: "Europe/Berlin", Active: true}, at(6, 12, 0), false},
		{"bad timezone", HoldWindow{Weekdays: weekdays, Start: "06:00", End: "07:30", Tz: "Nowhere/Land", Active: true}, at(2, 6, 45), false},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.w.ActiveAt(tc.ts))
		})
	}
}

func TestHoldWindowValidate(t *testing.T) {
	valid := HoldWindow{Weekdays: []int{1}, Start: "06:00", End: "07:00", Tz: "Europe/Berlin", MinCurrent: 6, MaxCurrent: 10}
	assert.NoError(t, valid.Validate())

	for name, mod := range map[string]func(*HoldWindow){
		"weekday":       func(w *HoldWindow) { w.Weekdays = []int{7} },
		"timezone":      func(w *HoldWindow) { w.Tz = "Nowhere/Land" },
		"start":         func(w *HoldWindow) { w.Start = "6am" },
		"end":           func(w *HoldWindow) { w.End = "25:00" },
		"equal":         func(w *HoldWindow) { w.End = w.Start },
		"negative":      func(w *HoldWindow) { w.MinCurrent = -1 },
		"min above max": func(w *HoldWindow) { w.MinCurrent, w.MaxCurrent = 10, 6 },
	} {
		t.Run(name, func(t *testing.T) {
			w := valid
			mod(&w)
			assert.Error(t, w.Validate())
		})
	}
}
