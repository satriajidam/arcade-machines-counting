package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type PlayTime struct {
	Begin time.Time
	End   time.Time
}

func (p1 PlayTime) String() string {
	bHour, bMinute, _ := p1.Begin.Clock()
	eHour, eMinute, _ := p1.End.Clock()
	return fmt.Sprintf("{%d:%02d - %d:%02d}", bHour, bMinute, eHour, eMinute)
}

func (p1 PlayTime) IsOverlapping(p2 PlayTime) bool {
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

func countArcadeMachines(playTimes []PlayTime) int {
	if len(playTimes) == 0 {
		return 0
	}

	overlappedCount := []int{}
	for i := 0; i < len(playTimes); i++ {
		count := 1

		for j := 0; j < len(playTimes); j++ {
			if i == j {
				continue
			}

			isOverlapped := playTimes[i].IsOverlapping(playTimes[j])

			if !isOverlapped && playTimes[i].End.Before(playTimes[i].Begin) {
				break
			}

			if isOverlapped {
				count++
			}
		}

		overlappedCount = append(overlappedCount, count)
	}

	highestCount := 0
	for idx, val := range overlappedCount {
		if idx == 0 || val > highestCount {
			highestCount = val
		}
	}

	return highestCount
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Error: no input & output file path provided")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: no input file path provided")
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: no provide output file provided")
		os.Exit(1)
	}

	inputFilepath := os.Args[1]
	inputFile, err := os.Open(inputFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed reading input file: %s\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	playTimes := []PlayTime{}
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		items := strings.Split(scanner.Text(), " ")
		pt, err := toPlayTime(items[0], items[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed reading input file: %s\n", err)
			os.Exit(1)
		}
		playTimes = append(playTimes, pt)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed reading input file: %s\n", err)
		os.Exit(1)
	}

	outputFilepath := os.Args[2]
	if _, err := os.Stat(outputFilepath); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: failed writing output file: create %s: file already exists\n", outputFilepath)
		os.Exit(1)
	}

	outputFile, err := os.OpenFile(outputFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed writing output file: %s\n", err)
		os.Exit(1)
	}

	writer := bufio.NewWriter(outputFile)
	_, err = writer.WriteString(fmt.Sprintf("%d\n", countArcadeMachines(playTimes)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed writing to output file: %s\n", err)
		os.Exit(1)
	}

	writer.Flush()
	outputFile.Close()
}
