package main

import (
	"testing"
	"time"
)

func TestToHHMM(t *testing.T) {
	tests := map[string]string{
		"":     "",
		"900":  "09:00",
		"910":  "09:10",
		"940":  "09:40",
		"1200": "12:00",
		"950":  "09:50",
		"1120": "11:20",
		"1100": "11:00",
		"1130": "11:30",
		"1300": "13:00",
		"1400": "14:00",
		"1350": "13:50",
		"1420": "14:20",
	}

	for rawHour, want := range tests {
		got := toHHMM(rawHour)
		if want != got {
			t.Errorf("toFormattedHour(\"%s\") got: %v, want: %v", rawHour, got, want)
		}
	}
}

func TestIsValidHHMM(t *testing.T) {
	invalidTests := []string{
		":", "99:10", "33:61", "100:00", "090:50", "011:200",
		"11:87", "11:130", "42:56", "140:00", "013:50", "014:200",
	}

	for _, hhmm := range invalidTests {
		got := isValidHHMM(hhmm)
		if got {
			t.Errorf("isValidHHMM(\"%s\") got: %v, want: %v", hhmm, got, false)
		}
	}

	validTests := []string{
		"09:00", "09:10", "09:40", "12:00", "09:50", "11:20",
		"11:00", "11:30", "13:00", "14:00", "13:50", "14:20",
	}

	for _, hhmm := range validTests {
		got := isValidHHMM(hhmm)
		if !got {
			t.Errorf("isValidHHMM(\"%s\") got: %v, want: %v", hhmm, got, true)
		}
	}
}

func TestToPlayTime(t *testing.T) {
	tests := []struct {
		name  string
		begin string
		end   string
		want  PlayTime
	}{
		{
			name:  "case1",
			begin: "900",
			end:   "910",
			want: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:10:00Z")
				return time
			}()},
		},
		{
			name:  "case2",
			begin: "940",
			end:   "1200",
			want: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:40:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T12:00:00Z")
				return time
			}()},
		},
		{
			name:  "case3",
			begin: "950",
			end:   "1120",
			want: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:50:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:20:00Z")
				return time
			}()},
		},
		{
			name:  "case4",
			begin: "1100",
			end:   "1130",
			want: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:30:00Z")
				return time
			}()},
		},
		{
			name:  "case5",
			begin: "1300",
			end:   "1400",
			want: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T13:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T14:00:00Z")
				return time
			}()},
		},
		{
			name:  "case6",
			begin: "1350",
			end:   "1420",
			want: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T13:50:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T14:20:00Z")
				return time
			}()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := toPlayTime(tt.begin, tt.end)
			if !got.Begin.Equal(tt.want.Begin) {
				t.Errorf("toPlayTime(\"%s\", \"%s\") got: %v, want: %v", tt.begin, tt.end, got.Begin, tt.want.Begin)
			}
			if !got.End.Equal(tt.want.End) {
				t.Errorf("toPlayTime(\"%s\", \"%s\") got: %v, want: %v", tt.begin, tt.end, got.End, tt.want.End)
			}
		})
	}
}

func TestIsOverlapping(t *testing.T) {
	tests := []struct {
		name string
		p1   PlayTime
		p2   PlayTime
		want bool
	}{
		{
			name: "09:00-09:10 is not overlapping with 09:40-12:00",
			p1: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:10:00Z")
				return time
			}()},
			p2: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:40:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T12:00:00Z")
				return time
			}()},
			want: false,
		},
		{
			name: "11:00-11:30 is not overlapping with 13:00-14:00",
			p1: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:30:00Z")
				return time
			}()},
			p2: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T13:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T14:00:00Z")
				return time
			}()},
			want: false,
		},
		{
			name: "09:40-12:00 is overlapping with 09:50-11:20",
			p1: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:40:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T12:00:00Z")
				return time
			}()},
			p2: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:50:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:20:00Z")
				return time
			}()},
			want: true,
		},
		{
			name: "09:40-12:00 is overlapping with 11:00-11:30",
			p1: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:40:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T12:00:00Z")
				return time
			}()},
			p2: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:30:00Z")
				return time
			}()},
			want: true,
		},
		{
			name: "09:50-11:20 is overlapping with 11:00-11:30",
			p1: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T09:50:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:20:00Z")
				return time
			}()},
			p2: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T11:30:00Z")
				return time
			}()},
			want: true,
		},
		{
			name: "13:00-14:00 is overlapping with 13:50-14:20",
			p1: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T13:00:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T14:00:00Z")
				return time
			}()},
			p2: PlayTime{func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T13:50:00Z")
				return time
			}(), func() time.Time {
				time, _ := time.Parse(time.RFC3339, "2006-01-03T14:20:00Z")
				return time
			}()},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p1.isOverlapping(tt.p2)
			if got != tt.want {
				t.Errorf("p1.isOverlapping(%v) got: %v, want: %v", tt.p2, got, tt.want)
			}
		})
	}
}