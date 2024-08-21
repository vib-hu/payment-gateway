package enums

import "fmt"

type PaymentMethod int

const (
	PaymentMethodCard      PaymentMethod = 1
	PaymentMethodApplePay  PaymentMethod = 2
	PaymentMethodGooglePay PaymentMethod = 3
)

func ParsePaymentMethod(value int) (PaymentMethod, error) {
	switch value {
	case int(PaymentMethodCard):
		return PaymentMethodCard, nil
	case int(PaymentMethodApplePay):
		return PaymentMethodApplePay, nil
	case int(PaymentMethodGooglePay):
		return PaymentMethodGooglePay, nil
	default:
		return 0, fmt.Errorf("invalid payment method: %d", value)
	}
}
