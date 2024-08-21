package commondtos

import value_objects "PaymentGateway/internal/payment-gateway-api/domains/value-objects"

type ApplePayDetailsRequest struct {
	ApplePayToken string `json:"apple_pay_token"`
}

func (applePay *ApplePayDetailsRequest) ToDomain() value_objects.ApplePay {
	if applePay == nil {
		return value_objects.ApplePay{}
	}

	return value_objects.ApplePay{
		Token: applePay.ApplePayToken,
	}
}
