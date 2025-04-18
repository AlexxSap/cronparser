package cronparser

import (
	"testing"
	"time"
)

func date(year int, month time.Month, day, hour, min int) time.Time {
	return time.Date(year, month, day, hour, min, 0, 0, time.UTC)
}

func TestIsMatch(t *testing.T) {
	tests := []struct {
		cronString string
		date       time.Time
		expected   bool
	}{
		// every minute is match
		{"* * * * *", date(2025, 4, 1, 12, 33), true},
		{"* * * * *", date(2012, 12, 12, 12, 12), true},

		// every 5 minutes
		{"/5 * * * *", date(2012, 12, 12, 12, 25), true},
		{"/5 * * * *", date(2012, 12, 12, 12, 26), false},

		//* 9 * * SAT - every minute, between 9:00 and 9:59, on saturday
		{"* 9 * * SAT", date(2025, 4, 19, 9, 12), true},
		{"* 9 * * SAT", date(2025, 4, 19, 10, 12), false},
		{"* 9 * * SAT", date(2025, 4, 18, 9, 12), false},

		//5 14-15 * 1,3 Mon-Fri - 14:05 and 15:05, on januar and march, from monday to friday
		{"5 14-15 * 1,3 mon-Fri", date(2025, 1, 16, 14, 5), true},
		{"5 14-15 * 1,3 mon-Fri", date(2025, 3, 14, 14, 5), true},
		{"5 14-15 * 1,3 mon-Fri", date(2025, 3, 15, 15, 1), false},
		{"5 14-15 * 1,3 mon-Fri", date(2025, 3, 14, 15, 6), false},

		///5 14 6-11 1 mon - every 5 minutes between 14:00 and 14:55, on 6-11 januar in monday
		{"/5 14 6-11 1 mon", date(2025, 1, 6, 14, 15), true},
		{"/5 14 6-11 1 mon", date(2025, 1, 6, 14, 35), true},
		{"/5 14 6-11 1 mon", date(2025, 1, 7, 14, 15), false},
		{"/5 14 6-11 1 mon", date(2025, 1, 6, 14, 13), false},
	}

	for _, test := range tests {
		/// TODO добавить проверку ошибок
		// actual, err := NewCronParser(test.cronString).IsMatch(test.date)
		// if err != nil {
		// t.Errorf("error: %v on string '%v' with date '%v'", err, test.cronString, test.date)
		// }

		actual, _ := NewCronParser(test.cronString).IsMatch(test.date)

		if actual != test.expected {
			t.Errorf("fail on string '%v' with date '%v'", test.cronString, test.date)
		}
	}
}

/// TODO список ошибок:
/// - неверное число токенов
/// - неверное семантическое значение для токена
/// - неверный символ в токене
/// - передана невалидная дата
