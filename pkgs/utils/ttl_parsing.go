package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	formatError    = `ttl string: invalid format "%s", ttl examples: "1h", "2h 45m", "30m"`
	unhandledError = `ttl string: unhandled error: %s`
)

var re = regexp.MustCompile(`^\s*(?P<hours>\d+h)?\s*(?P<minutes>\d+m)?\s*$`)

func TryParseTTL(ttl string) (time.Duration, error) {
	match := re.FindStringSubmatch(ttl)
	matchesCount := len(match)
	if match == nil || matchesCount != 3 || match[0] == "" {
		err := fmt.Errorf(formatError, ttl)
		return 0, err
	}
	return getValidResult(match[1], match[2])
}

func getValidResult(hours string, minutes string) (result time.Duration, err error) {
	var nHours, nMinutes time.Duration
	nHours, err = getTimePart(hours, "h")
	nMinutes, err = getTimePart(minutes, "m")
	return nHours*time.Hour + nMinutes*time.Minute, nil
}

func getTimePart(timePart string, timePartSuffix string) (time.Duration, error) {
	var result int
	var err error
	if timePart != "" {
		result, err = strconv.Atoi(strings.TrimSuffix(timePart, timePartSuffix))
		if err != nil {
			return 0, fmt.Errorf(unhandledError, err)
		}
	}
	return time.Duration(result), nil
}
