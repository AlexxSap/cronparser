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

	InvalidStepSize         = "InvalidStepSize"
	InvalidFromToValues     = "InvalidFromToSize"
	InvalidFromToValuesSize = "InvalidFromToValuesSize"
	InvalidListValues       = "InvalidListValues"
	InvalidSimpleValues     = "InvalidSimpleValues"

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

type valueType uint8

const (
	typeAll valueType = iota
	typeSteps
	typeFromTo
	typeLst
	typeSimple
)

type tokenType uint8

const (
	minutes tokenType = iota
	hour
	dayOfMonth
	month
	dayOfWeek
)

func valueOfDate(date time.Time, tType tokenType) uint8 {
	switch tType {
	case minutes:
		return uint8(date.Minute())
	case hour:
		return uint8(date.Hour())
	case dayOfMonth:
		return uint8(date.Day())
	case month:
		return uint8(date.Month())
	case dayOfWeek:
		return uint8(date.Weekday())
	}
	return 0
}

type token struct {
	vType     valueType
	minValue  uint8
	maxValue  uint8
	step      uint8
	lstValues []uint8
}

func newAllToken(minValue, maxValue uint8) (token, error) {
	return token{
			vType:    typeAll,
			minValue: minValue,
			maxValue: maxValue},
		nil
}

func newStepToken(str string, minValue, maxValue uint8) (token, error) {
	r, index := utf8.DecodeRuneInString(str)
	if r == utf8.RuneError {
		return token{}, errors.New(ErrorValidationUnknown)
	}
	stepInt, err := strconv.Atoi(str[index:])
	if err != nil {
		/// TODO добавить тест
		return token{}, fmt.Errorf("%v with: %w", InvalidStepSize, err)
	}

	step := uint8(stepInt)
	if step < minValue || step > maxValue {
		/// TODO добавить тест
		return token{}, errors.New(InvalidStepSize)
	}

	return token{
			vType:    typeSteps,
			step:     step,
			minValue: minValue,
			maxValue: maxValue},
		nil
}

func newFromToToken(str string, minValue, maxValue uint8) (token, error) {
	nums := strings.Split(str, symFromTo)
	if len(nums) != 2 {
		/// TODO добавить тест
		return token{}, errors.New(InvalidFromToValuesSize)
	}

	getValue := func(str string) (uint8, error) {
		valueInt, err := strconv.Atoi(str)
		if err != nil {
			/// TODO добавить тест
			return 0, fmt.Errorf("%v with: %w", InvalidFromToValues, err)
		}

		val := uint8(valueInt)
		if val < minValue || val > maxValue {
			/// TODO добавить тест
			return 0, errors.New(InvalidFromToValues)
		}
		return val, nil
	}

	valFrom, err := getValue(nums[0])
	if err != nil {
		return token{}, err
	}
	valTo, err := getValue(nums[1])
	if err != nil {
		/// TODO добавить тест
		return token{}, err
	}

	return token{
			vType:    typeFromTo,
			step:     0,
			minValue: valFrom,
			maxValue: valTo},
		nil
}

func newLstToken(str string, minValue, maxValue uint8) (token, error) {
	nums := strings.Split(str, symLst)

	values := make([]uint8, 0, len(nums))
	for _, num := range nums {
		valueInt, err := strconv.Atoi(num)
		if err != nil {
			/// TODO добавить тест
			return token{}, fmt.Errorf("%v with: %w", InvalidListValues, err)
		}

		val := uint8(valueInt)
		if val < minValue || val > maxValue {
			/// TODO добавить тест
			return token{}, errors.New(InvalidListValues)
		}
		values = append(values, val)
	}

	return token{
			vType:     typeLst,
			step:      0,
			minValue:  minValue,
			maxValue:  maxValue,
			lstValues: values},
		nil
}

func newSimpleNumberToken(str string, minValue, maxValue uint8) (token, error) {
	valueInt, err := strconv.Atoi(str)
	if err != nil {
		/// TODO добавить тест
		return token{}, fmt.Errorf("%v with: %w", InvalidSimpleValues, err)
	}

	val := uint8(valueInt)
	if val < minValue || val > maxValue {
		/// TODO добавить тест
		return token{}, errors.New(InvalidSimpleValues)
	}

	return token{
			vType:    typeSimple,
			step:     val,
			minValue: minValue,
			maxValue: maxValue},
		nil
}

func newToken(str string, minValue, maxValue uint8) (token, error) {
	if str == symAll {
		return newAllToken(minValue, maxValue)
	}

	if strings.HasPrefix(str, symStep) {
		return newStepToken(str, minValue, maxValue)
	}

	if strings.Contains(str, symFromTo) {
		return newFromToToken(str, minValue, maxValue)
	}

	if strings.Contains(str, symLst) {
		return newLstToken(str, minValue, maxValue)
	}

	return newSimpleNumberToken(str, minValue, maxValue)
}

func (tkn token) isMatch(dateTimeValue uint8) bool {
	switch tkn.vType {
	case typeAll:
		return true
	case typeSteps:
		return dateTimeValue%tkn.step == 0
	case typeFromTo:
		return false /// TODO дописать
	case typeLst:
		return false /// TODO дописать
	case typeSimple:
		return tkn.step == dateTimeValue
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

// IsMatch checks whether the specified date matches the expression.
// Parameters:
// - dateTime: object of time.Time to check
// Returns:
// - the result of checking (true or false)
// - error with failure case (if first result is false)
// Examples:
// res, err := NewCronParser("5 14-15 * 1,3 mon-Fri").IsMatch(time.Date(2025, 3, 14, 14, 5, 0, 0, time.UTC)) // result will be (true, nil)
// res, err := NewCronParser("5 14-15 * 1,3 mon-Fri").IsMatch(time.Date(2025, 3, 15, 15, 5, 0, 0, time.UTC)) // result will be (false, DayOfWeekNotMatch)
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
