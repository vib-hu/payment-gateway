package depositdtos

import (
	commondtos "PaymentGateway/internal/payment-gateway-api/controllers/models/common-dtos"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"log/slog"
)

type DepositRequest struct {
	Customer             commondtos.CustomerRequest    `json:"customer"`
	Amount               commondtos.AmountRequest      `json:"amount"`
	Country              enums.Country                 `json:"country_iso_code"`
	TransactionRouteType enums.TransactionRouteType    `json:"transaction_route_type"`
	BankDetails          commondtos.BankDetailsRequest `json:"bank_details"`
	CardDetails          commondtos.CardDetailsRequest `json:"card_details"`
	ClientReferenceId    string                        `json:"client_reference_id"`
}

func (p *DepositRequest) CountryIsUSA() bool {
	return p.Country == enums.CountryUSA
}

func (p *DepositRequest) CountryIsUK() bool {
	return p.Country == enums.CountryUK
}

// LogValue removes sensitive fields from log
func (p DepositRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("Amount", p.Amount),
		slog.Any("Country", p.Country),
		slog.Any("TransactionRouteType", p.TransactionRouteType),
		slog.Any("ClientReferenceId", p.ClientReferenceId),
	)
}
