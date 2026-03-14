package logic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

var ErrNoRepeat = errors.New("no repeat rule")

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return dateOnly.After(nowOnly)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrNoRepeat
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("Invalid date format: %w", err)
	}

	parts := strings.Fields(strings.TrimSpace(repeat))
	if len(parts) == 0 {
		return "", ErrNoRepeat
	}

	rule := parts[0]

	switch rule {
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid rule")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat rule")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid interval value")
		}

		if interval <= 0 || interval >= 400 {
			return "", fmt.Errorf("interval must be between 0 and 400")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("invalid repeat rule")
	}

	return date.Format(DateFormat), nil
}
