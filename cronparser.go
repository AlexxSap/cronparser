package cronparser

import "time"

type token struct {
}

type result struct {
	minute     token
	hour       token
	dayOfMonth token
	month      token
	dayOfWeek  token
}

type CronParser struct {
	cronString          string
	alreadyParsedResult result
	hasResult           bool
}

func NewCronParser(str string) CronParser {
	return CronParser{
		cronString:          str,
		alreadyParsedResult: result{},
		hasResult:           false}
}

func (crn CronParser) parse() (result, error) {
	if crn.hasResult {
		return crn.alreadyParsedResult, nil
	}

	return result{}, nil
}

func (res result) isMatch(dateTime time.Time) (bool, error) {
	return false, nil
}

func (crn CronParser) IsMatch(dateTime time.Time) (bool, error) {
	res, err := crn.parse()
	if err != nil {
		return false, err
	}

	isMatch, err := res.isMatch(dateTime)
	if err != nil {
		return false, err
	}

	return isMatch, nil
}
