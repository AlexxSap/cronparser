package cronparser

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	ErrorValidationEmptyString = "ErrorValidationEmptyString"
	ErrorValidationTokensCount = "ErrorValidationTokensCount"

	InvalidStepFormat       = "InvalidStepFormat"
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
	valueAll valueType = iota
	valueSteps
	valueFromTo
	valueLst
	valueSimple
)

type tokenType uint8

const (
	tokenMinutes tokenType = iota
	tokenHour
	tokenDayOfMonth
	tokenMonth
	tokenDayOfWeek
)

func valueOfDate(date time.Time, tType tokenType) uint8 {
	switch tType {
	case tokenMinutes:
		return uint8(date.Minute())
	case tokenHour:
		return uint8(date.Hour())
	case tokenDayOfMonth:
		return uint8(date.Day())
	case tokenMonth:
		return uint8(date.Month())
	case tokenDayOfWeek:
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

func nextSimpleValue(val uint8, simpleValue uint8) (uint8, bool) {
	return simpleValue, val > simpleValue
}

func nextAllValue(val, minVal, maxVal uint8) (uint8, bool) {
	if val < minVal {
		return minVal, false
	} else if val < maxVal {
		return val + 1, false
	} else {
		return minVal, true
	}
}

func nextStepValue(val, minVal, maxVal, step uint8) (uint8, bool) {
	if val < minVal {
		return minVal, false
	} else if val < maxVal {
		return ((val - minVal) / step) * step, false
	} else {
		return minVal, true
	}
}

func nextFromToValue(val, minVal, maxVal uint8) (uint8, bool) {
	if val < minVal {
		return minVal, false
	} else if val < maxVal {
		return val + 1, false
	} else {
		return minVal, true
	}
}

func nextLstValue(val uint8, lst []uint8) (uint8, bool) {
	if val < lst[0] {
		return lst[0], false
	} else if val < lst[len(lst)-1] {
		for _, lstVal := range lst {
			if val < lstVal {
				return lstVal, false
			}
		}
		return lst[0], true
	} else {
		return lst[0], true
	}
}

func (tkn token) next(currentValue uint8) (uint8, bool, uint8) {
	var nextValue uint8
	var loop bool
	minValue := tkn.minValue
	switch tkn.vType {
	case valueAll:
		nextValue, loop = nextAllValue(currentValue, tkn.minValue, tkn.maxValue)
	case valueSteps:
		nextValue, loop = nextStepValue(currentValue, tkn.minValue, tkn.maxValue, tkn.step)
	case valueFromTo:
		nextValue, loop = nextFromToValue(currentValue, tkn.minValue, tkn.maxValue)
	case valueLst:
		nextValue, loop = nextLstValue(currentValue, tkn.lstValues)
	case valueSimple:
		nextValue, loop = nextSimpleValue(currentValue, tkn.step)
		minValue = tkn.step
	}

	return nextValue, loop, minValue
}

func newAllToken(minValue, maxValue uint8) (token, error) {
	return token{
			vType:    valueAll,
			minValue: minValue,
			maxValue: maxValue},
		nil
}

func newStepToken(str string, minValue, maxValue uint8) (token, error) {
	r, index := utf8.DecodeRuneInString(str)
	if r == utf8.RuneError {
		return token{}, errors.New(InvalidStepFormat)
	}
	stepInt, err := strconv.Atoi(str[index:])
	if err != nil {
		return token{}, fmt.Errorf("%v with: %w", InvalidStepSize, err)
	}

	step := uint8(stepInt)
	if step < minValue || step > maxValue {
		return token{}, errors.New(InvalidStepSize)
	}

	return token{
			vType:    valueSteps,
			step:     step,
			minValue: minValue,
			maxValue: maxValue},
		nil
}

func newFromToToken(str string, minValue, maxValue uint8) (token, error) {
	nums := strings.Split(str, symFromTo)
	if len(nums) != 2 {
		return token{}, errors.New(InvalidFromToValuesSize)
	}

	getValue := func(str string) (uint8, error) {
		valueInt, err := strconv.Atoi(str)
		if err != nil {
			return 0, fmt.Errorf("%v with: %w", InvalidFromToValues, err)
		}

		val := uint8(valueInt)
		if val < minValue || val > maxValue {
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
		return token{}, err
	}

	if valTo <= valFrom {
		return token{}, errors.New(InvalidFromToValues)
	}

	return token{
			vType:    valueFromTo,
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
			return token{}, fmt.Errorf("%v with: %w", InvalidListValues, err)
		}

		val := uint8(valueInt)
		if val < minValue || val > maxValue {
			return token{}, errors.New(InvalidListValues)
		}
		values = append(values, val)
	}
	slices.Sort(values)

	return token{
			vType:     valueLst,
			step:      0,
			minValue:  minValue,
			maxValue:  maxValue,
			lstValues: values},
		nil
}

func newSimpleNumberToken(str string, minValue, maxValue uint8) (token, error) {
	valueInt, err := strconv.Atoi(str)
	if err != nil {
		return token{}, fmt.Errorf("%v with: %w", InvalidSimpleValues, err)
	}

	val := uint8(valueInt)
	if val < minValue || val > maxValue {
		return token{}, errors.New(InvalidSimpleValues)
	}

	return token{
			vType:    valueSimple,
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
	case valueAll:
		return true
	case valueSteps:
		return dateTimeValue%tkn.step == 0
	case valueFromTo:
		return dateTimeValue >= tkn.minValue && dateTimeValue <= tkn.maxValue
	case valueLst:
		return slices.Contains(tkn.lstValues, dateTimeValue)
	case valueSimple:
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

func replaceMonthNamesToNumber(str string) string {
	lstr := strings.ToLower(str)
	lstr = strings.ReplaceAll(lstr, "jan", "1")
	lstr = strings.ReplaceAll(lstr, "feb", "2")
	lstr = strings.ReplaceAll(lstr, "mar", "3")
	lstr = strings.ReplaceAll(lstr, "apr", "4")
	lstr = strings.ReplaceAll(lstr, "may", "5")
	lstr = strings.ReplaceAll(lstr, "jun", "6")
	lstr = strings.ReplaceAll(lstr, "jul", "7")
	lstr = strings.ReplaceAll(lstr, "aug", "8")
	lstr = strings.ReplaceAll(lstr, "sep", "9")
	lstr = strings.ReplaceAll(lstr, "oct", "10")
	lstr = strings.ReplaceAll(lstr, "nov", "11")
	lstr = strings.ReplaceAll(lstr, "dec", "12")

	return lstr
}

func replaceDayNamesToNumber(str string) string {
	lstr := strings.ToLower(str)
	lstr = strings.ReplaceAll(lstr, "mon", "1")
	lstr = strings.ReplaceAll(lstr, "tue", "2")
	lstr = strings.ReplaceAll(lstr, "wed", "3")
	lstr = strings.ReplaceAll(lstr, "thu", "4")
	lstr = strings.ReplaceAll(lstr, "fri", "5")
	lstr = strings.ReplaceAll(lstr, "sat", "6")
	lstr = strings.ReplaceAll(lstr, "sun", "7")
	return lstr
}

func tokenize(str string) (result, error) {
	if len(str) == 0 {
		return result{}, errors.New(ErrorValidationEmptyString)
	}

	strTokens := strings.Fields(str)
	if len(strTokens) != 5 {
		return result{}, errors.New(ErrorValidationTokensCount)
	}

	minute, err := newToken(strTokens[0], 0, 59)
	if err != nil {
		return result{}, err
	}

	hour, err := newToken(strTokens[1], 0, 23)
	if err != nil {
		return result{}, err
	}

	dayOfMonth, err := newToken(strTokens[2], 1, 31)
	if err != nil {
		return result{}, err
	}

	month, err := newToken(replaceMonthNamesToNumber(strTokens[3]), 1, 12)
	if err != nil {
		return result{}, err
	}

	dayOfWeek, err := newToken(replaceDayNamesToNumber(strTokens[4]), 1, 7)
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
	if !res.minute.isMatch(valueOfDate(dateTime, tokenMinutes)) {
		return false, errors.New(MinutesNotMatch)
	}

	if !res.hour.isMatch(valueOfDate(dateTime, tokenHour)) {
		return false, errors.New(HoursNotMatch)
	}

	if !res.dayOfMonth.isMatch(valueOfDate(dateTime, tokenDayOfMonth)) {
		return false, errors.New(DayOfMonthNotMatch)
	}

	if !res.month.isMatch(valueOfDate(dateTime, tokenMonth)) {
		return false, errors.New(MonthNotMatch)
	}

	if !res.dayOfWeek.isMatch(valueOfDate(dateTime, tokenDayOfWeek)) {
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

// TODO добавить комментарии
// TODO добавить исполняемый пример
// TODO добавить бенчмарк на сравнение с простым перебором дат
func (crn CronParser) NearestDate(currentDate time.Time) (time.Time, error) {
	parsed, err := crn.parse()
	if err != nil {
		return time.Time{}, err
	}

	/// TODO добавить проверку - нужно ли проверять следующий токен
	newMinute, loop, minMinute := parsed.minute.next(valueOfDate(currentDate, tokenMinutes))
	if !loop {
		/// TODO придумать способ нормально задавать дату
		newDate := time.Date(
			currentDate.Year(),
			currentDate.Month(),
			currentDate.Day(),
			currentDate.Hour(),
			int(newMinute),
			0,
			0,
			currentDate.Location())
		if res, _ := parsed.isMatch(newDate); res {
			return newDate, nil
		}
	}
	newHour, loop, minHour := parsed.hour.next(valueOfDate(currentDate, tokenHour))
	if loop {
		newMinute = minMinute
	}

	newDay, loop, minDay := parsed.dayOfMonth.next(valueOfDate(currentDate, tokenDayOfMonth))
	if loop {
		newHour = minHour
		newMinute = minMinute
	}

	newMonth, loop, _ := parsed.month.next(valueOfDate(currentDate, tokenMonth))
	/// TODO добавить проверку дня недели
	// for !parsed.dayOfWeek.isMatch(newDay) {
	// newDay, loop, minDay = parsed.dayOfMonth.next()
	// }

	if loop {
		newDay = minDay
		newHour = minHour
		newMinute = minMinute
	}

	newYear := currentDate.Year()
	if loop {
		newYear = newYear + 1
	}

	fmt.Println("year", newYear)
	fmt.Println("month", newMonth)
	fmt.Println("day", newDay)
	fmt.Println("hour", newHour)
	fmt.Println("min", newMinute)

	return time.Date(newYear, time.Month(newMonth), int(newDay), int(newHour), int(newMinute), 0, 0, currentDate.Location()), nil
}
