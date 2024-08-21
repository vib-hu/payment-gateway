package enums

import "fmt"

type Currency int

const (
	CurrencyUSA     Currency = 840
	CurrencyCanada  Currency = 124
	CurrencyUK      Currency = 826
	CurrencyGermany Currency = 276
	CurrencyIndia   Currency = 356
)

func ParseCurrency(value int) (Currency, error) {
	switch value {
	case int(CurrencyUSA):
		return CurrencyUSA, nil
	case int(CurrencyCanada):
		return CurrencyCanada, nil
	case int(CurrencyUK):
		return CurrencyUK, nil
	case int(CurrencyGermany):
		return CurrencyGermany, nil
	case int(CurrencyIndia):
		return CurrencyIndia, nil
	default:
		return 0, fmt.Errorf("invalid currency: %d", value)
	}
}
