package cronparser

import (
	"testing"
	"time"
)

func TestIsMatch(t *testing.T) {
	tests := []struct {
		cronString string
		date       time.Time
		expected   bool
	}{
		{"* * * * *", time.Date(2025, 4, 1, 12, 33, 1, 0, time.UTC), true},
		{"* * * * *", time.Date(2012, 12, 12, 12, 12, 12, 12, time.UTC), true},
	}

	for _, test := range tests {
		actual, err := NewCronParser(test.cronString).IsMatch(test.date)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		if actual != test.expected {
			t.Errorf("actual (%v) != expected (%v)", actual, test.expected)
		}
	}
}

/// TODO список ошибок:
/// - неверное число токенов
/// - неверное семантическое значение для токена
/// - неверный символ в токене
/// - передана невалидная дата
