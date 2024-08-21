package valueobjects

import "PaymentGateway/internal/payment-gateway-api/enums"

type Amount struct {
	Value    float64
	Currency enums.Currency
}

func (amount *Amount) IsInvalid() bool {
	return amount == nil || amount.Value < 0
}
