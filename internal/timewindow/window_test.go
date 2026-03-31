package timewindow

import (
	"testing"
	"time"
)

func TestInBlockedWindow_SameDay(t *testing.T) {
	loc := time.UTC

	tests := []struct {
		name        string
		now         time.Time
		window      string
		wantBlocked bool
		wantAfter   time.Duration
	}{
		{
			name:        "before window",
			now:         time.Date(2026, 3, 19, 12, 0, 0, 0, loc),
			window:      "13:00-15:00",
			wantBlocked: false,
		},
		{
			name:        "inside window",
			now:         time.Date(2026, 3, 19, 14, 0, 0, 0, loc),
			window:      "13:00-15:00",
			wantBlocked: true,
			wantAfter:   time.Hour,
		},
		{
			name:        "at start boundary",
			now:         time.Date(2026, 3, 19, 13, 0, 0, 0, loc),
			window:      "13:00-15:00",
			wantBlocked: true,
			wantAfter:   2 * time.Hour,
		},
		{
			name:        "at end boundary",
			now:         time.Date(2026, 3, 19, 15, 0, 0, 0, loc),
			window:      "13:00-15:00",
			wantBlocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocked, after, err := InBlockedWindow(tt.now, tt.window)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if blocked != tt.wantBlocked {
				t.Fatalf("blocked=%v want=%v", blocked, tt.wantBlocked)
			}
			if blocked && after != tt.wantAfter {
				t.Fatalf("after=%v want=%v", after, tt.wantAfter)
			}
		})
	}
}

func TestInBlockedWindow_Overnight(t *testing.T) {
	loc := time.UTC

	tests := []struct {
		name        string
		now         time.Time
		window      string
		wantBlocked bool
		wantAfter   time.Duration
	}{
		{
			name:        "before overnight start",
			now:         time.Date(2026, 3, 19, 22, 0, 0, 0, loc),
			window:      "23:00-02:00",
			wantBlocked: false,
		},
		{
			name:        "inside overnight before midnight",
			now:         time.Date(2026, 3, 19, 23, 30, 0, 0, loc),
			window:      "23:00-02:00",
			wantBlocked: true,
			wantAfter:   2*time.Hour + 30*time.Minute,
		},
		{
			name:        "inside overnight after midnight",
			now:         time.Date(2026, 3, 19, 1, 0, 0, 0, loc),
			window:      "23:00-02:00",
			wantBlocked: true,
			wantAfter:   time.Hour,
		},
		{
			name:        "after overnight end",
			now:         time.Date(2026, 3, 19, 3, 0, 0, 0, loc),
			window:      "23:00-02:00",
			wantBlocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocked, after, err := InBlockedWindow(tt.now, tt.window)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if blocked != tt.wantBlocked {
				t.Fatalf("blocked=%v want=%v", blocked, tt.wantBlocked)
			}
			if blocked && after != tt.wantAfter {
				t.Fatalf("after=%v want=%v", after, tt.wantAfter)
			}
		})
	}
}

func TestInBlockedWindow_Invalid(t *testing.T) {
	_, _, err := InBlockedWindow(time.Now(), "bad-value")
	if err == nil {
		t.Fatal("expected error")
	}
}
