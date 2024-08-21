package models

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"fmt"
)

type GatewayBWithdrawRequest struct {
	TransactionID  string  `xml:"transaction_id"`
	Amount         float64 `xml:"amount"`
	Currency       int     `xml:"currency"`
	PaymentMethod  string  `xml:"payment_method"`
	Card           Card    `xml:"card_number,omitempty"`
	ApplePayToken  string  `xml:"apple_pay_token"`
	GooglePayToken string  `xml:"google_pay_token"`
	BillingAddress string  `xml:"billing_address"`
}

func ConvertToGatewayBWithdrawRequest(withdraw *domains.Withdraw) (*GatewayBWithdrawRequest, error) {
	if withdraw == nil {
		return nil, fmt.Errorf("withdraw cannot be nil")
	}

	pspRequest := &GatewayBWithdrawRequest{
		TransactionID:  withdraw.Id.String(),
		Amount:         withdraw.Amount.Value,
		Currency:       int(withdraw.Amount.Currency),
		BillingAddress: withdraw.BillingAddress,
		PaymentMethod:  string(withdraw.PaymentMethod),
	}

	switch withdraw.PaymentMethod {
	case enums.PaymentMethodCard:
		pspRequest.Card = Card{
			CardNumber: withdraw.Card.CardNumber,
			Cvv:        withdraw.Card.Cvv,
		}
		break
	case enums.PaymentMethodApplePay:
		pspRequest.ApplePayToken = withdraw.ApplePay.Token
		break
	case enums.PaymentMethodGooglePay:
		pspRequest.GooglePayToken = withdraw.GooglePay.Token
		break
	}

	return pspRequest, nil
}
