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

type Deposit struct {
	Id                     uuid.UUID
	Amount                 value_objects.Amount
	Country                enums.Country
	TransactionRouteType   enums.TransactionRouteType
	Customer               Customer
	Card                   value_objects.Card
	BankAccount            value_objects.BankAccount
	TransactionStatus      enums.TransactionStatus
	ClientReferenceId      string
	GatewayIdentifier      string
	PspResponseCode        string
	PspResponseDescription string
	Description            string
	CreatedDateUtc         time.Time
	ModifiedDateUtc        time.Time
}

func NewDeposit(id uuid.UUID, amount value_objects.Amount, country enums.Country, routeType enums.TransactionRouteType,
	customer Customer, card value_objects.Card, bankAccount value_objects.BankAccount,
	gatewayIdentifier string, clientReferenceId string) (*Deposit, error) {

	err := depositDomainValidations(amount, customer, routeType, bankAccount, card, clientReferenceId, gatewayIdentifier)
	if err != nil {
		return nil, err
	}

	return &Deposit{
		Id:                   id,
		Amount:               amount,
		Country:              country,
		TransactionRouteType: routeType,
		Customer:             customer,
		Card:                 card,
		BankAccount:          bankAccount,
		GatewayIdentifier:    gatewayIdentifier,
		ClientReferenceId:    clientReferenceId,
		TransactionStatus:    enums.TransactionInitiated,
		CreatedDateUtc:       time.Now().UTC(),
		ModifiedDateUtc:      time.Now().UTC(),
	}, nil
}

func (d *Deposit) Initiate() error {
	if d.TransactionStatus != enums.TransactionInitiated {
		return errors.New("invalid transaction status")
	}

	d.TransactionStatus = enums.TransactionInitiated
	d.ModifiedDateUtc = time.Now().UTC()
	return nil
}

func (d *Deposit) Success(pspResponseCode string, pspResponseDescription string) error {
	err := d.transactionUpdateValidations(pspResponseCode, pspResponseDescription)
	if err != nil {
		return err
	}

	d.TransactionStatus = enums.TransactionSuccessful
	d.updatePspDetails(pspResponseCode, pspResponseDescription)
	d.ModifiedDateUtc = time.Now().UTC()
	return nil
}

func (d *Deposit) FailedDueToPayment(pspResponseCode string, pspResponseDescription string) error {
	err := d.transactionUpdateValidations(pspResponseCode, pspResponseDescription)
	if err != nil {
		return err
	}

	d.TransactionStatus = enums.TransactionFailed
	d.updatePspDetails(pspResponseCode, pspResponseDescription)
	d.ModifiedDateUtc = time.Now().UTC()
	return nil
}

func (d *Deposit) FailedDueToInternalError(description string) error {
	if d.TransactionStatus != enums.TransactionInitiated {
		return errors.New("invalid transaction status")
	}
	d.TransactionStatus = enums.TransactionFailed
	d.Description = description
	d.ModifiedDateUtc = time.Now().UTC()
	return nil
}

func (d *Deposit) UpdateStatus(status enums.TransactionStatus) {
	d.TransactionStatus = status
	d.ModifiedDateUtc = time.Now().UTC()
}

func depositDomainValidations(amount value_objects.Amount, customer Customer, routeType enums.TransactionRouteType,
	bankAccount value_objects.BankAccount, card value_objects.Card, clientReferenceId string, gatewayIdentifier string) error {
	if amount.IsInvalid() {
		return errors.New("amount details are invalid")
	}
	if customer.IsInvalid() {
		return errors.New("customer details are invalid")
	}
	if routeType != enums.RouteToBank && routeType != enums.RouteToCard {
		return errors.New("invalid transaction route type")
	}

	if routeType == enums.RouteToBank && bankAccount.IsInvalid() {
		return errors.New("invalid bank account details")
	}

	if routeType == enums.RouteToCard && card.IsInvalid() {
		return errors.New("invalid card details")
	}

	if strings.TrimSpace(clientReferenceId) == "" {
		return errors.New("invalid client reference id")
	}

	if strings.TrimSpace(gatewayIdentifier) == "" {
		return errors.New("invalid gateway identifier")
	}
	return nil
}

func (d *Deposit) transactionUpdateValidations(pspResponseCode string, pspResponseDescription string) error {
	if d.TransactionStatus != enums.TransactionInitiated {
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

func (d *Deposit) updatePspDetails(pspResponseCode string, pspResponseDescription string) {
	d.PspResponseCode = pspResponseCode
	d.PspResponseDescription = pspResponseDescription
}

// LogValue removes sensitive fields from log
func (d Deposit) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("Id", d.Id),
		slog.Any("Amount", d.Amount.Value),
		slog.Any("Currency", d.Amount.Currency),
		slog.Any("TransactionRouteType", d.TransactionRouteType),
		slog.Any("TransactionStatus", d.TransactionStatus),
	)
}
