package helpers

import (
	"fmt"
)

// NumericStringToTiyin конвертирует строку вида "12345.67" или "12345" в тийины (int64)
func NumericStringToTiyin(s string) (int64, error) {
	var sign int64 = 1
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		s = s[1:]
	}

	var whole int64
	var frac int64
	if n, _ := fmt.Sscanf(s, "%d.%d", &whole, &frac); n == 2 {
		if frac < 0 || frac > 99 {
			return 0, fmt.Errorf("invalid fractional part: %s", s)
		}
		return sign * (whole*100 + frac), nil
	}
	if n, _ := fmt.Sscanf(s, "%d", &whole); n == 1 {
		return sign * (whole * 100), nil
	}

	return 0, fmt.Errorf("invalid numeric format: %s", s)
}

// TiyinToNumericString конвертирует тийины (int64) в строку "XXX.YY"
func TiyinToNumericString(t int64) string {
	sign := ""
	if t < 0 {
		sign = "-"
		t = -t
	}
	return fmt.Sprintf("%s%d.%02d", sign, t/100, t%100)
}
