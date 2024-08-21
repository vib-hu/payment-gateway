package models

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"fmt"
	"log/slog"
)

type GatewayADepositRequest struct {
	TransactionID string      `json:"transaction_id"`
	Amount        float64     `json:"amount"`
	Currency      int         `json:"currency"`
	Card          Card        `json:"card_number,omitempty"`
	BankAccount   BankAccount `json:"bank_account,omitempty"`
}

type Card struct {
	CardNumber string `json:"card_number"`
	Cvv        string `json:"cvv"`
}

type BankAccount struct {
	AccountNumber     string `json:"account_number"`
	RoutingNumber     string `json:"routing_number"`
	AccountHolderName string `json:"account_holder_name"`
}

func ConvertToGatewayADepositRequest(deposit *domains.Deposit) (*GatewayADepositRequest, error) {
	if deposit == nil {
		return nil, fmt.Errorf("deposit cannot be nil")
	}

	pspRequest := &GatewayADepositRequest{
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

// LogValue removes sensitive fields from log
func (p GatewayADepositRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("TransactionID", p.TransactionID),
		slog.Any("Amount", p.Amount),
		slog.Any("Currency", p.Currency),
	)
}
