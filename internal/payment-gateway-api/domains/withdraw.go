package domains

import (
	value_objects "PaymentGateway/internal/payment-gateway-api/domains/value-objects"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"strings"
	"time"
)

type Withdraw struct {
	Id                     uuid.UUID
	Amount                 value_objects.Amount
	Country                enums.Country
	PaymentMethod          enums.PaymentMethod
	Customer               Customer
	Card                   value_objects.Card
	ApplePay               value_objects.ApplePay
	GooglePay              value_objects.GooglePay
	TransactionStatus      enums.TransactionStatus
	BillingAddress         string
	ClientReferenceId      string
	GatewayIdentifier      string
	PspResponseCode        string
	PspResponseDescription string
	Description            string
	CreatedDateUtc         time.Time
	ModifiedDateUtc        time.Time
}

func NewWithdraw(id uuid.UUID, amount value_objects.Amount, country enums.Country, customer Customer,
	paymentMethod enums.PaymentMethod, card value_objects.Card, applePay value_objects.ApplePay, googlePay value_objects.GooglePay,
	gatewayIdentifier string, clientReferenceId string, billingAddress string) (*Withdraw, error) {

	err := withdrawDomainValidations(amount, customer, paymentMethod, card, applePay, googlePay, clientReferenceId,
		gatewayIdentifier, billingAddress)
	if err != nil {
		return nil, err
	}

	return &Withdraw{
		Id:                id,
		Amount:            amount,
		Country:           country,
		PaymentMethod:     paymentMethod,
		Customer:          customer,
		Card:              card,
		ApplePay:          applePay,
		GooglePay:         googlePay,
		GatewayIdentifier: gatewayIdentifier,
		ClientReferenceId: clientReferenceId,
		BillingAddress:    billingAddress,
		TransactionStatus: enums.TransactionInitiated,
		CreatedDateUtc:    time.Now().UTC(),
		ModifiedDateUtc:   time.Now().UTC(),
	}, nil
}

func (w *Withdraw) Initiate() error {
	if w.TransactionStatus != enums.TransactionInitiated {
		return errors.New("invalid transaction status")
	}

	w.TransactionStatus = enums.TransactionInitiated
	w.ModifiedDateUtc = time.Now().UTC()
	return nil
}

func (w *Withdraw) Success(pspResponseCode string, pspResponseDescription string) error {
	err := w.transactionUpdateValidations(pspResponseCode, pspResponseDescription)
	if err != nil {
		return err
	}

	w.TransactionStatus = enums.TransactionSuccessful
	w.updatePspDetails(pspResponseCode, pspResponseDescription)
	w.ModifiedDateUtc = time.Now().UTC()
	return nil
}

func (w *Withdraw) FailedDueToPayment(pspResponseCode string, pspResponseDescription string) error {
	err := w.transactionUpdateValidations(pspResponseCode, pspResponseDescription)
	if err != nil {
		return err
	}

	w.TransactionStatus = enums.TransactionFailed
	w.updatePspDetails(pspResponseCode, pspResponseDescription)
	w.ModifiedDateUtc = time.Now().UTC()
	return nil
}

func (w *Withdraw) FailedDueToInternalError(description string) error {
	if w.TransactionStatus != enums.TransactionInitiated {
		return errors.New("invalid transaction status")
	}
	w.TransactionStatus = enums.TransactionFailed
	w.Description = description
	w.ModifiedDateUtc = time.Now().UTC()
	return nil
}

func (d *Withdraw) UpdateStatus(status enums.TransactionStatus) {
	d.TransactionStatus = status
	d.ModifiedDateUtc = time.Now().UTC()
}

func withdrawDomainValidations(amount value_objects.Amount, customer Customer, paymentMethod enums.PaymentMethod,
	card value_objects.Card, applePay value_objects.ApplePay, googlePay value_objects.GooglePay, clientReferenceId string,
	gatewayIdentifier string, billingAddress string) error {
	if amount.IsInvalid() {
		return errors.New("amount details are invalid")
	}
	if customer.IsInvalid() {
		return errors.New("customer details are invalid")
	}
	if paymentMethod == enums.PaymentMethodCard && card.IsInvalid() {
		return errors.New("invalid card details")
	}
	if paymentMethod == enums.PaymentMethodApplePay && applePay.IsInvalid() {
		return errors.New("invalid apple pay details")
	}
	if paymentMethod == enums.PaymentMethodGooglePay && googlePay.IsInvalid() {
		return errors.New("invalid google pay details")
	}
	if strings.TrimSpace(clientReferenceId) == "" {
		return errors.New("invalid client reference id")
	}
	if strings.TrimSpace(billingAddress) == "" {
		return errors.New("invalid billing address")
	}
	if strings.TrimSpace(gatewayIdentifier) == "" {
		return errors.New("invalid gateway identifier")
	}
	return nil
}

func (w *Withdraw) transactionUpdateValidations(pspResponseCode string, pspResponseDescription string) error {
	if w.TransactionStatus != enums.TransactionInitiated {
		return errors.New("invalid transaction status")
	}

	if strings.TrimSpace(pspResponseCode) == "" {
		return errors.New("invalid psp response code")
	}

	if strings.TrimSpace(pspResponseDescription) == "" {
		return errors.New("invalid psp response description")
	}
	return nil
}

func (w *Withdraw) updatePspDetails(pspResponseCode string, pspResponseDescription string) {
	w.PspResponseCode = pspResponseCode
	w.PspResponseDescription = pspResponseDescription
}

// LogValue removes sensitive fields from log
func (w Withdraw) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("Id", w.Id),
		slog.Any("Amount", w.Amount.Value),
		slog.Any("Currency", w.Amount.Currency),
		slog.Any("PaymentMethod", w.PaymentMethod),
		slog.Any("TransactionStatus", w.TransactionStatus),
		slog.Any("GatewayIdentifier", w.GatewayIdentifier),
	)
}
