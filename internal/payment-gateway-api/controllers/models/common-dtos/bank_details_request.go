package commondtos

import value_objects "PaymentGateway/internal/payment-gateway-api/domains/value-objects"

type BankDetailsRequest struct {
	AccountNumber     string `json:"account_number"`
	RoutingNumber     string `json:"routing_number"`
	AccountHolderName string `json:"account_holder_name"`
}

func (bankDetails *BankDetailsRequest) ToDomain() value_objects.BankAccount {
	if bankDetails == nil {
		return value_objects.BankAccount{}
	}

	return value_objects.BankAccount{
		AccountNumber:     bankDetails.AccountNumber,
		RoutingNumber:     bankDetails.RoutingNumber,
		AccountHolderName: bankDetails.AccountHolderName,
	}
}
