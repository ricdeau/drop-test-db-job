package internal

import (
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

type Schedule struct {
	cron.Schedule
}

func (s *Schedule) UnmarshalJSON(bytes []byte) (err error) {
	spec := strings.Trim(string(bytes), `"`)
	s.Schedule, err = cron.ParseStandard(strings.ToLower(spec))
	if err != nil {
		return fmt.Errorf("parse cron schedule spec %q: %v", spec, err)
	}

	return nil
}
