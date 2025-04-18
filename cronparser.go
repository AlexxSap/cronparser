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

// type valueType int

const (
	typeAll = iota
	typeSteps
	typeFromTo
	typeLst
)

// type tokenType int
const (
	minutes = iota
	hour
	dayOfMonth
	month
	dayOfWeek
)

func valueOfDate(date time.Time, typeOfToken int) int {
	if typeOfToken == minutes {
		return date.Minute()
	}

	if typeOfToken == hour {
		return date.Hour()
	}

	if typeOfToken == dayOfMonth {
		return date.Day()
	}

	if typeOfToken == month {
		return int(date.Month())
	}

	if typeOfToken == dayOfWeek {
		return int(date.Weekday())
	}

	return 0
}

type token struct {
	tokenType int
	minValue  int
	maxValue  int
	step      int
}

func newToken(str string, minValue, maxValue int) (token, error) {
	if str == symAll {
		return token{
				tokenType: typeAll,
				minValue:  minValue,
				maxValue:  maxValue},
			nil
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

		if step < minValue || step > maxValue {
			/// TODO добавить тест
			return token{}, errors.New(InvalidStepSize)
		}

		return token{
				tokenType: typeSteps,
				step:      step,
				minValue:  minValue,
				maxValue:  maxValue},
			nil
	}

	if strings.Contains(str, symFromTo) {

	}

	if strings.Contains(str, symLst) {

	}

	/// TODO добавить тест
	return token{}, errors.New(ErrorValidationUnknown)
}

func (tkn token) isMatch(dateTimeValue int) bool {
	if tkn.tokenType == typeAll {
		return true
	}

	if tkn.tokenType == typeSteps {
		if dateTimeValue%tkn.step == 0 {
			return true
		}
	}

	if tkn.tokenType == typeFromTo {
		return false
	}

	if tkn.tokenType == typeLst {
		return false
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

	minute, err := newToken(strTokens[0], 1, 59)
	if err != nil {
		return result{}, err
	}

	hour, err := newToken(strTokens[1], 1, 23)
	if err != nil {
		return result{}, err
	}

	dayOfMonth, err := newToken(strTokens[2], 1, 31)
	if err != nil {
		return result{}, err
	}

	month, err := newToken(strTokens[3], 1, 12)
	if err != nil {
		return result{}, err
	}

	dayOfWeek, err := newToken(strTokens[4], 1, 7)
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
	if !res.minute.isMatch(valueOfDate(dateTime, minutes)) {
		/// TODO добавить тест
		return false, errors.New(MinutesNotMatch)
	}

	if !res.hour.isMatch(valueOfDate(dateTime, hour)) {
		/// TODO добавить тест
		return false, errors.New(HoursNotMatch)
	}

	if !res.dayOfMonth.isMatch(valueOfDate(dateTime, dayOfMonth)) {
		/// TODO добавить тест
		return false, errors.New(DayOfMonthNotMatch)
	}

	if !res.month.isMatch(valueOfDate(dateTime, month)) {
		/// TODO добавить тест
		return false, errors.New(MonthNotMatch)
	}

	if !res.dayOfWeek.isMatch(valueOfDate(dateTime, dayOfWeek)) {
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
