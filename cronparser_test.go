package cronparser

import (
	"fmt"
	"testing"
	"time"
)

func date(year int, month time.Month, day, hour, min int) time.Time {
	return time.Date(year, month, day, hour, min, 0, 0, time.UTC)
}

func errToStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func TestIsMatch(t *testing.T) {
	tests := []struct {
		cronString string
		date       time.Time
		expected   bool
		err        string
	}{
		// every minute is match
		{"* * * * *", date(2025, 4, 1, 12, 33), true, ""},
		{"* * * * *", date(2012, 12, 12, 12, 12), true, ""},

		// every 5 minutes
		{"/5 * * * *", date(2012, 12, 12, 12, 25), true, ""},
		{"/5 * * * *", date(2012, 12, 12, 12, 26), false, MinutesNotMatch},

		//* 9 * * SAT - every minute, between 9:00 and 9:59, on saturday
		{"* 9 * * SAT", date(2025, 4, 19, 9, 12), true, ""},
		{"* 9 * * SAT", date(2025, 4, 19, 10, 12), false, HoursNotMatch},
		{"* 9 * * SAT", date(2025, 4, 18, 9, 12), false, DayOfWeekNotMatch},
		{"* 9 19 * SAT", date(2025, 4, 18, 9, 12), false, DayOfMonthNotMatch},
		{"* 9 19 3 SAT", date(2025, 4, 19, 9, 12), false, MonthNotMatch},

		//5 14-15 * 1,3 Mon-Fri - 14:05 and 15:05, on januar and march, from monday to friday
		{"5 14-15 * 1,3 mon-Fri", date(2025, 1, 16, 14, 5), true, ""},
		{"5 14-15 * 1,3 mon-Fri", date(2025, 3, 14, 14, 5), true, ""},
		{"5 14-15 * 1,3 1-5", date(2025, 3, 14, 14, 5), true, ""},
		{"5 14-15 * 1,3 mon-Fri", date(2025, 3, 15, 15, 1), false, MinutesNotMatch},
		{"5 14-15 * 1,3 mon-Fri", date(2025, 3, 14, 15, 6), false, MinutesNotMatch},

		///5 14 6-11 1 mon - every 5 minutes between 14:00 and 14:55, on 6-11 januar in monday
		{"/5 14 6-11 1 mon", date(2025, 1, 6, 14, 15), true, ""},
		{"/5 14 6-11 1 mon", date(2025, 1, 6, 14, 35), true, ""},
		{"/5 14 6-11 1 mon", date(2025, 1, 7, 14, 15), false, DayOfWeekNotMatch},
		{"/5 14 6-11 1 mon", date(2025, 1, 6, 14, 13), false, MinutesNotMatch},

		/// * * 30 1,2 * - 30 jan or 30 feb. but 30 feb does not exists
		{"* * 30 1,2 *", date(2025, 1, 30, 0, 0), true, ""},
		{"* * 30 1,2 *", date(2025, 2, 30, 0, 0), false, DayOfMonthNotMatch},
	}

	for _, test := range tests {
		actual, err := NewCronParser(test.cronString).IsMatch(test.date)
		if actual != test.expected {
			t.Errorf("fail on string '%v' with date '%v' with err: %v", test.cronString, test.date, err)
		} else if errToStr(err) != test.err {
			t.Errorf("fail on string '%v' with date '%v' with err: %v (expected error: '%v')", test.cronString, test.date, err, test.err)
		}
	}
}

func TestValidationTokenize(t *testing.T) {
	tests := []struct {
		str string
		err string
	}{
		{"", ErrorValidationEmptyString},
		{"* * * *", ErrorValidationTokensCount},
		{"/5 * * 1-2 * *", ErrorValidationTokensCount},
		{"/a * * * *", "InvalidStepSize with: strconv.Atoi: parsing \"a\": invalid syntax"},
		{"/77 * * * *", InvalidStepSize},
		{"* * /-33 * *", InvalidStepSize},
		{"1-2-3 * * * *", InvalidFromToValuesSize},
		{"1- * * * *", "InvalidFromToSize with: strconv.Atoi: parsing \"\": invalid syntax"},
		{"1-a * * * *", "InvalidFromToSize with: strconv.Atoi: parsing \"a\": invalid syntax"},
		{"1-223 * * * *", InvalidFromToValues},
		{"111-22 * * * *", InvalidFromToValues},
		{"22-11 * * * *", InvalidFromToValues},
		{"1,a * * * *", "InvalidListValues with: strconv.Atoi: parsing \"a\": invalid syntax"},
		{"1,11,111 * * * *", InvalidListValues},
		{"1a * * * *", "InvalidSimpleValues with: strconv.Atoi: parsing \"1a\": invalid syntax"},
		{"111 * * * *", InvalidSimpleValues},
	}

	for _, test := range tests {
		_, err := tokenize(test.str)
		if err.Error() != test.err {
			t.Errorf("fail on string '%v' with error '%v' (expected '%v')", test.str, err, test.err)
		}
	}
}

func TestNearestDate(t *testing.T) {
	tests := []struct {
		str          string
		currentDate  time.Time
		expectedDate time.Time
		err          string
	}{
		// {"25 2,8,3 * * *", date(2025, 2, 6, 2, 12), date(2025, 2, 6, 2, 25), ""},
		{"25 2,8,3 * * *", date(2025, 2, 6, 1, 12), date(2025, 2, 6, 2, 25), ""},
		// {"25 2,8,3 * * *", date(2025, 2, 6, 1, 12), date(2025, 2, 6, 2, 25), ""},
		// {"25 2,8,3 * * *", date(2025, 2, 6, 2, 25), date(2025, 2, 6, 3, 25), ""},
		// {"/5 14 6-11 1 *", date(2025, 2, 6, 14, 12), date(2025, 1, 6, 14, 15), ""},
		// {"/5 14 6-11 1 *", date(2025, 1, 6, 13, 55), date(2025, 1, 6, 14, 0), ""},
		// {"/5 14 6-11 1 *", date(2025, 1, 6, 15, 55), date(2025, 1, 7, 14, 0), ""},
		// {"/5 14 6-11 1,4 *", date(2025, 1, 12, 15, 55), date(2025, 4, 6, 14, 0), ""},
		// {"/5 14 6-11 4,1 *", date(2025, 1, 12, 15, 55), date(2025, 4, 6, 14, 0), ""},
		// {"/5 14 * * mon,sat", date(2025, 4, 26, 15, 55), date(2025, 4, 28, 14, 0), ""},
	}

	for _, test := range tests {
		act, err := NewCronParser(test.str).NearestDate(test.currentDate)
		if act != test.expectedDate {
			t.Errorf("fail on string '%v' expected '%v' get '%v' with error '%v'", test.str, test.expectedDate, act, err)
		} else if errToStr(err) != test.err {
			t.Errorf("fail on string '%v' with error '%v' (expected '%v')", test.str, err, test.err)
		}
	}
}

func TestDateLoop(t *testing.T) {
	act := date(2025, 13, 12, 0, 0)
	exp := date(2026, 1, 12, 0, 0)

	if act != exp {
		t.Error("dates not match")
	}

	act = date(2025, 1, 12, 13, 77)
	exp = date(2025, 1, 12, 14, 17)

	if act != exp {
		t.Error("times not match")
	}
}

func BenchmarkIsMatch(b *testing.B) {
	pattern := "%d %d 1-%d 3-%d 2-5"
	for i := 0; i < b.N; i++ {
		NewCronParser(fmt.Sprintf(pattern, i%60, i%24, i%25, i%12)).IsMatch(date(2025, 1, 6, 14, 15))
	}
}

func ExampleCronParser_IsMatch() {
	actual, err := NewCronParser("/5 14 6-11 1 mon").IsMatch(time.Date(2025, 1, 6, 14, 35, 0, 0, time.UTC))
	fmt.Println("actual:", actual)
	fmt.Println("err:", err)
	// Output:
	// actual: true
	// err: <nil>
}

func ExampleCronParser_IsMatch_second() {
	actual, err := NewCronParser("/5 14 6-11 1 mon").IsMatch(time.Date(2025, 1, 7, 14, 35, 0, 0, time.UTC))
	fmt.Println("actual:", actual)
	fmt.Println("err:", err)
	// Output:
	// actual: false
	// err: DayOfWeekNotMatch
}
