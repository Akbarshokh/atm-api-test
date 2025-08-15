package domain

import (
	"errors"
	"fmt"
)

// деньги будут храниться в тийн
type Money int64

func (m Money) Tiyin() int64 { return int64(m) }

func (m Money) SumString() string {
	sign := ""
	v := m
	if v < 0 {
		sign = "-"
		v = -v
	}
	whole := int64(v) / 100
	cents := int64(v) % 100
	return fmt.Sprintf("%s%d.%02d", sign, whole, cents)
}

func FromSum(s string) (Money, error) {
	var whole int64
	var frac int64
	if _, err := fmt.Sscanf(s, "%d.%d", &whole, &frac); err == nil {
		if frac < 0 || frac > 99 {
			return 0, errors.New("invalid fractional part")
		}
		return Money(whole*100 + frac), nil
	}
	if _, err := fmt.Sscanf(s, "%d", &whole); err == nil {
		return Money(whole * 100), nil
	}
	return 0, errors.New("invalid sum format")
}
