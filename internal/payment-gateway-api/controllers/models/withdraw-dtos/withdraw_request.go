package withdrawdtos

import (
	commondtos "PaymentGateway/internal/payment-gateway-api/controllers/models/common-dtos"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"log/slog"
)

type WithdrawRequest struct {
	Customer          commondtos.CustomerRequest         `json:"customer"`
	Amount            commondtos.AmountRequest           `json:"amount"`
	Country           enums.Country                      `json:"country_iso_code"`
	PaymentMethod     enums.PaymentMethod                `json:"payment_method"`
	ApplePayDetails   commondtos.ApplePayDetailsRequest  `json:"apple_pay_details"`
	GooglePayDetails  commondtos.GooglePayDetailsRequest `json:"google_pay_details"`
	CardDetails       commondtos.CardDetailsRequest      `json:"card_details"`
	ClientReferenceId string                             `json:"client_reference_id"`
	BillingAddress    string                             `json:"billing_address"`
}

func (p *WithdrawRequest) CountryIsUSA() bool {
	return p.Country == enums.CountryUSA
}

func (p *WithdrawRequest) CountryIsUK() bool {
	return p.Country == enums.CountryUK
}

// LogValue removes sensitive fields from log
func (p WithdrawRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("Amount", p.Amount),
		slog.Any("Country", p.Country),
		slog.Any("PaymentMethod", p.PaymentMethod),
		slog.Any("ClientReferenceId", p.ClientReferenceId),
		slog.Any("BillingAddress", p.BillingAddress),
	)
}
