package scheduler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	if y1 > y2 {
		return true
	}
	if y1 == y2 && m1 > m2 {
		return true
	}
	if y1 == y2 && m1 == m2 && d1 > d2 {
		return true
	}

	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {

	case "d":

		if len(parts) != 2 {
			return "", errors.New("invalid repeat format")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if interval < 1 || interval > 400 {
			return "", errors.New("invalid day interval")
		}

		for {
			date = date.AddDate(0, 0, interval)

			if afterNow(date, now) {
				break
			}
		}

	case "y":

		for {
			date = date.AddDate(1, 0, 0)

			if afterNow(date, now) {
				break
			}
		}

	default:
		return "", errors.New("unsupported repeat rule")
	}

	return date.Format("20060102"), nil
}
