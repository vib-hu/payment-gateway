package commondtos

import value_objects "PaymentGateway/internal/payment-gateway-api/domains/value-objects"

type CardDetailsRequest struct {
	CardNumber string `json:"card_number"`
	Cvv        string `json:"cvv"`
}

func (cardDetails *CardDetailsRequest) ToDomain() value_objects.Card {
	if cardDetails == nil {
		return value_objects.Card{}
	}
	return value_objects.Card{
		CardNumber: cardDetails.CardNumber,
		Cvv:        cardDetails.Cvv,
	}
}
