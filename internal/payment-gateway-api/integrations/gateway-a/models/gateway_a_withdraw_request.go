package models

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"fmt"
)

type GatewayAWithdrawRequest struct {
	TransactionID  string  `json:"transaction_id"`
	Amount         float64 `json:"amount"`
	Currency       int     `json:"currency"`
	PaymentMethod  string  `json:"payment_method"`
	Card           Card    `json:"card_number,omitempty"`
	ApplePayToken  string  `json:"apple_pay_token"`
	GooglePayToken string  `json:"google_pay_token"`
	BillingAddress string  `json:"billing_address"`
}

func ConvertToGatewayAWithdrawRequest(withdraw *domains.Withdraw) (*GatewayAWithdrawRequest, error) {
	if withdraw == nil {
		return nil, fmt.Errorf("withdraw cannot be nil")
	}

	pspRequest := &GatewayAWithdrawRequest{
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
