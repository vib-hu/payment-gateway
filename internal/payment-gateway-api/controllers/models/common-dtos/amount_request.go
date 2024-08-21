package commondtos

import (
	value_objects "PaymentGateway/internal/payment-gateway-api/domains/value-objects"
	"PaymentGateway/internal/payment-gateway-api/enums"
)

type AmountRequest struct {
	Value    float64        `json:"value"`
	Currency enums.Currency `json:"currency_iso_code"`
}

func (amountRequest *AmountRequest) ToDomain() value_objects.Amount {
	if amountRequest == nil {
		return value_objects.Amount{}
	}

	return value_objects.Amount{
		Value:    amountRequest.Value,
		Currency: amountRequest.Currency,
	}
}
