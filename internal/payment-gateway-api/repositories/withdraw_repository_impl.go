package repositories

import (
	"PaymentGateway/internal/payment-gateway-api/constants"
	"PaymentGateway/internal/payment-gateway-api/domains"
	value_objects "PaymentGateway/internal/payment-gateway-api/domains/value-objects"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"PaymentGateway/pkg/encryption"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"log/slog"
	"os"
	"strings"
)

type WithdrawRepositoryImpl struct {
	db        *sql.DB
	encryptor encryption.Encryption
	logger    *slog.Logger
}

func NewWithdrawRepository(encryptor encryption.Encryption, logger *slog.Logger) (*WithdrawRepositoryImpl, error) {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if strings.TrimSpace(connectionString) == "" {
		errorMsg := "database connection string is missing"
		logger.Error(errorMsg)
		return nil, errors.New(errorMsg)
	}

	// verifying if db source exists using DB.Ping without creating a connection
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		errorMsg := fmt.Sprintf("error while connecting to the database: %s", err.Error())
		logger.Error(errorMsg)
		return nil, err
	}

	return &WithdrawRepositoryImpl{db: db, encryptor: encryptor}, nil
}

func (w *WithdrawRepositoryImpl) Create(ctx context.Context, withdraw *domains.Withdraw) error {
	fieldsToEncrypt := map[string]string{
		constants.GooglePayToken:  withdraw.GooglePay.Token,
		constants.ApplePayToken:   withdraw.ApplePay.Token,
		constants.CardNumberField: withdraw.Card.CardNumber,
		constants.CvvField:        withdraw.Card.Cvv,
		constants.BillingAddress:  withdraw.BillingAddress,
	}
	encryptedFields, err := w.encryptor.EncryptMultiple(fieldsToEncrypt)
	if err != nil {
		return err
	}
	query := `
        INSERT INTO withdraw (
            id, amount, currency, country, customer_id, payment_method, card_number, cvv, apple_pay_token, google_pay_token,
            status, client_reference_id, gateway_identifier, psp_response_code, psp_response_description, description,
            billing_address, created_date_utc, modified_date_utc
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	_, err = w.db.ExecContext(ctx, query, withdraw.Id, withdraw.Amount.Value, withdraw.Amount.Currency, withdraw.Country,
		withdraw.Customer.Id, withdraw.PaymentMethod, encryptedFields[constants.CardNumberField], encryptedFields[constants.CvvField],
		encryptedFields[constants.ApplePayToken], encryptedFields[constants.GooglePayToken], withdraw.TransactionStatus,
		withdraw.ClientReferenceId, withdraw.GatewayIdentifier, withdraw.PspResponseCode, withdraw.PspResponseDescription,
		withdraw.Description, encryptedFields[constants.BillingAddress], withdraw.CreatedDateUtc, withdraw.ModifiedDateUtc)
	return err
}

func (w *WithdrawRepositoryImpl) Update(ctx context.Context, withdraw *domains.Withdraw) error {
	query := `
       UPDATE withdraw SET status=$2, psp_response_code=$3, psp_response_description=$4, description=$5, modified_date_utc = $6
       WHERE id = $1`
	_, err := w.db.ExecContext(ctx, query, withdraw.Id, withdraw.TransactionStatus, withdraw.PspResponseCode, withdraw.PspResponseDescription,
		withdraw.Description, withdraw.ModifiedDateUtc)
	return err
}

func (w *WithdrawRepositoryImpl) GetById(ctx context.Context, id uuid.UUID) (*domains.Withdraw, error) {
	query := `
        SELECT 
            id, amount, currency, country, customer_id, payment_method, card_number, cvv, apple_pay_token, google_pay_token, 
            status, client_reference_id, gateway_identifier, psp_response_code, psp_response_description, description,
            billing_address, created_date_utc, modified_date_utc
        FROM 
            withdraw
        WHERE 
            id = $1
    `

	var withdraw domains.Withdraw
	var amountValue float64
	var amountCurrency string
	var customerId int64
	var currency, country, transactionStatus, paymentMethod int
	var cardNumber, cardCvv, applePayToken, googlePayToken, billingAddress string

	err := w.db.QueryRow(query, id).Scan(
		&withdraw.Id, &amountValue, &amountCurrency, &country, &customerId, &paymentMethod,
		&cardNumber, &cardCvv, &applePayToken, &googlePayToken, &transactionStatus,
		&withdraw.ClientReferenceId, &withdraw.GatewayIdentifier,
		&withdraw.PspResponseCode, &withdraw.PspResponseDescription, &withdraw.Description,
		&billingAddress, &withdraw.CreatedDateUtc, &withdraw.ModifiedDateUtc)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No withdraw found with the given ID
		}
		return nil, err
	}
	parsedCurrency, err := enums.ParseCurrency(currency)
	withdraw.Amount = value_objects.Amount{
		Value:    amountValue,
		Currency: parsedCurrency,
	}
	withdraw.Country, _ = enums.ParseCountry(country)
	withdraw.TransactionStatus, _ = enums.ParseTransactionStatus(transactionStatus)
	withdraw.Customer = domains.Customer{Id: customerId}

	decryptedCardNumber, _ := w.encryptor.Decrypt(cardNumber)
	decryptedCvv, _ := w.encryptor.Decrypt(cardCvv)
	decryptedApplePayToken, _ := w.encryptor.Decrypt(applePayToken)
	decryptedGooglePayToken, _ := w.encryptor.Decrypt(googlePayToken)
	decryptedBillingAddress, _ := w.encryptor.Decrypt(billingAddress)
	withdraw.PaymentMethod = enums.PaymentMethod(paymentMethod)
	withdraw.Card = value_objects.Card{
		CardNumber: decryptedCardNumber,
		Cvv:        decryptedCvv,
	}
	withdraw.ApplePay = value_objects.ApplePay{
		Token: decryptedApplePayToken,
	}
	withdraw.GooglePay = value_objects.GooglePay{
		Token: decryptedGooglePayToken,
	}
	withdraw.BillingAddress = decryptedBillingAddress

	return &withdraw, nil
}
