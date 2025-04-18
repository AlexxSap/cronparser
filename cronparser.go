package cronparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	ErrorValidationEmptyString = "ErrorValidationEmptyString"
	ErrorValidationTokensCount = "ErrorValidationTokensCount"
	ErrorValidationUnknown     = "ErrorValidationUnknown"

	InvalidStepSize = "InvalidStepSize"

	MinutesNotMatch    = "MinutesNotMatch"
	HoursNotMatch      = "HoursNotMatch"
	DayOfMonthNotMatch = "DayOfMonthNotMatch"
	MonthNotMatch      = "MonthNotMatch"
	DayOfWeekNotMatch  = "DayOfWeekNotMatch"
)

const (
	symAll    = "*"
	symStep   = "/"
	symFromTo = "-"
	symLst    = ","
)

const (
	typeAll = iota
	typeSteps
	typeFromTo
	typeLst
)

type token struct {
	tokenType int
	step      int
}

func newToken(str string) (token, error) {
	if str == symAll {
		return token{tokenType: typeAll}, nil
	}

	if strings.HasPrefix(str, symStep) {
		r, index := utf8.DecodeRuneInString(str)
		if r == utf8.RuneError {
			return token{}, errors.New(ErrorValidationUnknown)
		}
		step, err := strconv.Atoi(str[index:])
		if err != nil {
			/// TODO добавить тест
			return token{}, fmt.Errorf("%v with: %w", InvalidStepSize, err)
		}

		if step <= 0 || step > 59 {
			/// TODO добавить тест
			return token{}, errors.New(InvalidStepSize)
		}

		return token{tokenType: typeSteps, step: step}, nil
	}

	if strings.Contains(str, symFromTo) {

	}

	if strings.Contains(str, symLst) {

	}

	/// TODO добавить тест
	return token{}, errors.New(ErrorValidationUnknown)
}

func (tkn token) isMatch(dateTime time.Time) bool {
	if tkn.tokenType == typeAll {
		return true
	}

	if tkn.tokenType == typeSteps {
		if dateTime.Minute()%tkn.step == 0 {
			return true
		}
	}

	return false
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
	if !res.minute.isMatch(dateTime) {
		/// TODO добавить тест
		return false, errors.New(MinutesNotMatch)
	}

	if !res.hour.isMatch(dateTime) {
		/// TODO добавить тест
		return false, errors.New(HoursNotMatch)
	}

	if !res.dayOfMonth.isMatch(dateTime) {
		/// TODO добавить тест
		return false, errors.New(DayOfMonthNotMatch)
	}

	if !res.month.isMatch(dateTime) {
		/// TODO добавить тест
		return false, errors.New(MonthNotMatch)
	}

	if !res.dayOfWeek.isMatch(dateTime) {
		/// TODO добавить тест
		return false, errors.New(DayOfWeekNotMatch)
	}

	return true, nil
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
