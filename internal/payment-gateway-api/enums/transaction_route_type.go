package enums

import "fmt"

type TransactionRouteType int

const (
	RouteToBank TransactionRouteType = 1
	RouteToCard TransactionRouteType = 2
)

func ParseTransactionRouteType(value int) (TransactionRouteType, error) {
	switch value {
	case int(RouteToBank):
		return RouteToBank, nil
	case int(RouteToCard):
		return RouteToCard, nil
	default:
		return 0, fmt.Errorf("invalid transaction route type: %d", value)
	}
}
