package cronparser

import (
	"errors"
	"strings"
	"time"
)

const (
	ErrorValidationEmptyString = "ErrorValidationEmptyString"
	ErrorValidationTokensCount = "ErrorValidationTokensCount"
)

type token struct {
}

func newToken(str string) (token, error) {
	return token{}, nil
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

func tokenize(str string) (result, error) {
	if len(str) == 0 {
		/// TODO добавить тест
		return result{}, errors.New(ErrorValidationEmptyString)
	}

	strTokens := strings.Fields(str)
	if len(strTokens) != 5 {
		/// TODO добавить тест
		return result{}, errors.New(ErrorValidationTokensCount)
	}

	minute, err := newToken(strTokens[0])
	if err != nil {
		return result{}, err
	}

	hour, err := newToken(strTokens[1])
	if err != nil {
		return result{}, err
	}

	dayOfMonth, err := newToken(strTokens[2])
	if err != nil {
		return result{}, err
	}

	month, err := newToken(strTokens[3])
	if err != nil {
		return result{}, err
	}

	dayOfWeek, err := newToken(strTokens[4])
	if err != nil {
		return result{}, err
	}

	return result{
		minute:     minute,
		hour:       hour,
		dayOfMonth: dayOfMonth,
		month:      month,
		dayOfWeek:  dayOfWeek,
	}, nil
}

func (crn *CronParser) parse() (result, error) {
	if crn.hasResult {
		return crn.alreadyParsedResult, nil
	}

	parsed, err := tokenize(crn.cronString)
	if err != nil {
		return result{}, err
	}

	crn.alreadyParsedResult = parsed
	crn.hasResult = true

	return parsed, nil
}

func (res result) isMatch(dateTime time.Time) (bool, error) {
	return false, nil
}

func (crn CronParser) IsMatch(dateTime time.Time) (bool, error) {
	parsed, err := crn.parse()
	if err != nil {
		return false, err
	}

	isMatch, err := parsed.isMatch(dateTime)
	if err != nil {
		return false, err
	}

	return isMatch, nil
}
