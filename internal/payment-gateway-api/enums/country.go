package enums

import "fmt"

type Country int

// Country ISO codes
const (
	CountryUSA     Country = 840
	CountryCanada  Country = 124
	CountryUK      Country = 826
	CountryGermany Country = 276
	CountryIndia   Country = 356
)

func ParseCountry(value int) (Country, error) {
	switch value {
	case int(CountryUSA):
		return CountryUSA, nil
	case int(CountryCanada):
		return CountryCanada, nil
	case int(CountryUK):
		return CountryUK, nil
	case int(CountryGermany):
		return CountryGermany, nil
	case int(CountryIndia):
		return CountryIndia, nil
	default:
		return 0, fmt.Errorf("invalid country: %d", value)
	}
}
