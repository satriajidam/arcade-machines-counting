package main

import (
	"fmt"
	"regexp"
	"time"
)

type PlayTime struct {
	Begin time.Time
	End   time.Time
}

func (p1 PlayTime) isOverlapping(p2 PlayTime) bool {
	return p1.Begin.Before(p2.End) && p2.Begin.Before(p1.End)
}

func toHHMM(rawHour string) string {
	if len(rawHour) == 3 {
		return fmt.Sprintf("0%s:%s", string(rawHour[0]), rawHour[1:])
	}

	if len(rawHour) == 4 {
		return fmt.Sprintf("%s:%s", rawHour[:2], rawHour[2:])
	}

	return ""
}

func isValidHHMM(hhmm string) bool {
	matched, _ := regexp.MatchString(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`, hhmm)
	return matched
}

func toPlayTime(beginRaw, endRaw string) (PlayTime, error) {
	beginHHMM := toHHMM(beginRaw)
	if !isValidHHMM(beginHHMM) {
		return PlayTime{}, fmt.Errorf("Invalid hh:mm format for beginning time: %s", beginHHMM)
	}

	endHHMM := toHHMM(endRaw)
	if !isValidHHMM(endHHMM) {
		return PlayTime{}, fmt.Errorf("Invalid hh:mm format for ending time: %s", endHHMM)
	}

	beginTime, err := time.Parse(time.RFC3339, fmt.Sprintf("2006-01-03T%s:00Z", beginHHMM))
	if err != nil {
		return PlayTime{}, err
	}

	endTime, err := time.Parse(time.RFC3339, fmt.Sprintf("2006-01-03T%s:00Z", endHHMM))
	if err != nil {
		return PlayTime{}, err
	}

	return PlayTime{beginTime, endTime}, nil
}

func main() {

}
