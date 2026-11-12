package api

import (
	"fmt"
	"slices"
	"time"
)

// HoldWindow is a recurring time window in which the charger is kept enabled
// (e.g. to power the vehicle's own scheduled pre-heating) without charging
// towards a soc target. Min/MaxCurrent of 0 mean the loadpoint's effective values.
type HoldWindow struct {
	Weekdays   []int   `json:"weekdays"`   // 0-6 (Sunday-Saturday), weekday the window starts on
	Start      string  `json:"start"`      // HH:MM
	End        string  `json:"end"`        // HH:MM, an end before start spans midnight
	Tz         string  `json:"tz"`         // timezone in IANA format
	MinCurrent float64 `json:"minCurrent"` // lower current limit in A, 0 for default
	MaxCurrent float64 `json:"maxCurrent"` // upper current limit in A, 0 for default
	Active     bool    `json:"active"`     // active flag
}

// Validate checks the window for consistency
func (w HoldWindow) Validate() error {
	for _, day := range w.Weekdays {
		if day < 0 || day > 6 {
			return fmt.Errorf("weekday out of range: %v", day)
		}
	}
	if _, err := time.LoadLocation(w.Tz); err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}
	start, err := time.Parse("15:04", w.Start)
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}
	end, err := time.Parse("15:04", w.End)
	if err != nil {
		return fmt.Errorf("invalid end time: %w", err)
	}
	if start.Equal(end) {
		return fmt.Errorf("start and end must differ")
	}
	if w.MinCurrent < 0 || w.MaxCurrent < 0 {
		return fmt.Errorf("current must not be negative")
	}
	if w.MaxCurrent > 0 && w.MinCurrent > w.MaxCurrent {
		return fmt.Errorf("min current %.3gA exceeds max current %.3gA", w.MinCurrent, w.MaxCurrent)
	}
	return nil
}

// ActiveAt returns whether the window is active at the given time.
// The start is inclusive, the end exclusive. Windows ending before they start span midnight.
func (w HoldWindow) ActiveAt(ts time.Time) bool {
	if !w.Active || len(w.Weekdays) == 0 {
		return false
	}

	loc, err := time.LoadLocation(w.Tz)
	if err != nil {
		return false
	}
	start, err := time.Parse("15:04", w.Start)
	if err != nil {
		return false
	}
	end, err := time.Parse("15:04", w.End)
	if err != nil {
		return false
	}

	ts = ts.In(loc)
	minutes := ts.Hour()*60 + ts.Minute()
	startMin := start.Hour()*60 + start.Minute()
	endMin := end.Hour()*60 + end.Minute()

	switch {
	case startMin < endMin:
		return slices.Contains(w.Weekdays, int(ts.Weekday())) && minutes >= startMin && minutes < endMin
	case startMin > endMin:
		// first part runs on the start weekday, remainder after midnight belongs to the previous weekday
		if minutes >= startMin {
			return slices.Contains(w.Weekdays, int(ts.Weekday()))
		}
		if minutes < endMin {
			return slices.Contains(w.Weekdays, int(ts.AddDate(0, 0, -1).Weekday()))
		}
	}

	return false
}
