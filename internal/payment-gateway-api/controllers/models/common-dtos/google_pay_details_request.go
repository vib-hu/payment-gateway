package commondtos

import value_objects "PaymentGateway/internal/payment-gateway-api/domains/value-objects"

type GooglePayDetailsRequest struct {
	GooglePayToken string `json:"google_pay_token"`
}

func (googlePay *GooglePayDetailsRequest) ToDomain() value_objects.GooglePay {
	if googlePay == nil {
		return value_objects.GooglePay{}
	}

	return value_objects.GooglePay{
		Token: googlePay.GooglePayToken,
	}
}
