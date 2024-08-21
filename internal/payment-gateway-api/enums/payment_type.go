package enums

import (
	"errors"
	"strings"
)

type PaymentType string

const (
	Deposit  PaymentType = "deposit"
	Withdraw PaymentType = "withdraw"
)

func ParseToPaymentType(paymentType string) (PaymentType, error) {
	switch strings.ToLower(paymentType) {
	case string(Deposit):
		return Deposit, nil
	case string(Withdraw):
		return Withdraw, nil
	default:
		return "", errors.New("invalid payment type")
	}
}
