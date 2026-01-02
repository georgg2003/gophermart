package models

import "math"

type Money struct {
	AmountMinor int64
}

func (m Money) Major() float64 {
	return float64(m.AmountMinor) / 100
}

func NewMoneyFromMajor(r float64) Money {
	return Money{AmountMinor: int64(math.Round(r * 100))}
}

func NewMoney(m int64) Money {
	return Money{AmountMinor: m}
}
