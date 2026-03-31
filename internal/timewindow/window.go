package timewindow

import (
	"fmt"
	"time"
)

func InBlockedWindow(now time.Time, value string) (bool, time.Duration, error) {
	startTOD, endTOD, err := parseDailyWindow(value)
	if err != nil {
		return false, 0, err
	}

	y, m, d := now.Date()
	loc := now.Location()

	startToday := time.Date(y, m, d, startTOD.Hour(), startTOD.Minute(), 0, 0, loc)
	endToday := time.Date(y, m, d, endTOD.Hour(), endTOD.Minute(), 0, 0, loc)

	// Same-day window, e.g. 13:00-15:00
	if endToday.After(startToday) {
		if (now.Equal(startToday) || now.After(startToday)) && now.Before(endToday) {
			return true, endToday.Sub(now), nil
		}
		return false, 0, nil
	}

	// Overnight window, e.g. 23:00-02:00
	endNextDay := endToday.Add(24 * time.Hour)
	if now.Equal(startToday) || now.After(startToday) {
		return true, endNextDay.Sub(now), nil
	}

	startYesterday := startToday.Add(-24 * time.Hour)
	if now.After(startYesterday) && now.Before(endToday) {
		return true, endToday.Sub(now), nil
	}

	return false, 0, nil
}

func parseDailyWindow(value string) (time.Time, time.Time, error) {
	const layout = "15:04"

	var startStr, endStr string
	n, err := fmt.Sscanf(value, "%5s-%5s", &startStr, &endStr)
	if err != nil || n != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid window format %q, expected HH:MM-HH:MM", value)
	}

	startTOD, err := time.Parse(layout, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time %q", startStr)
	}

	endTOD, err := time.Parse(layout, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time %q", endStr)
	}

	return startTOD, endTOD, nil
}
