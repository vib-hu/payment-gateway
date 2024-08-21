package models

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"fmt"
)

type GatewayBDepositRequest struct {
	TransactionID string      `xml:"transaction_id"`
	Amount        float64     `xml:"amount"`
	Currency      int         `xml:"currency"`
	Card          Card        `xml:"card_number,omitempty"`
	BankAccount   BankAccount `xml:"bank_account,omitempty"`
}

type Card struct {
	CardNumber string `xml:"card_number"`
	Cvv        string `xml:"cvv"`
}

type BankAccount struct {
	AccountNumber     string `xml:"account_number"`
	RoutingNumber     string `xml:"routing_number"`
	AccountHolderName string `xml:"account_holder_name"`
}

func ConvertToGatewayBDepositRequest(deposit *domains.Deposit) (*GatewayBDepositRequest, error) {
	if deposit == nil {
		return nil, fmt.Errorf("deposit domain should be present")
	}

	pspRequest := &GatewayBDepositRequest{
		TransactionID: deposit.Id.String(),
		Amount:        deposit.Amount.Value,
		Currency:      int(deposit.Amount.Currency),
	}

	if deposit.TransactionRouteType == enums.RouteToCard {
		pspRequest.Card = Card{
			CardNumber: deposit.Card.CardNumber,
			Cvv:        deposit.Card.Cvv,
		}
	} else if deposit.TransactionRouteType == enums.RouteToBank {
		pspRequest.BankAccount = BankAccount{
			AccountNumber:     deposit.BankAccount.AccountNumber,
			RoutingNumber:     deposit.BankAccount.RoutingNumber,
			AccountHolderName: deposit.BankAccount.AccountHolderName,
		}
	}

	return pspRequest, nil
}
