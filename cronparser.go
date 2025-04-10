package cronparser

import "time"

type CronParser struct {
	cronString string
}

func (crn CronParser) IsMatch(time.Time) (bool, error) {
	return true, nil
}
