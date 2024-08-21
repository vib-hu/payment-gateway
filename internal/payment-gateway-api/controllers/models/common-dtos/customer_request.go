package commondtos

import (
	value_objects "PaymentGateway/internal/payment-gateway-api/domains"
)

type CustomerRequest struct {
	Id int64
}

func (customerRequest *CustomerRequest) ToDomain() value_objects.Customer {
	if customerRequest == nil {
		return value_objects.Customer{}
	}

	return value_objects.Customer{
		Id: customerRequest.Id,
	}
}
