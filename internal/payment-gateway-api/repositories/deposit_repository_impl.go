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

type DepositRepositoryImpl struct {
	db        *sql.DB
	encryptor encryption.Encryption
	logger    *slog.Logger
}

func NewDepositRepository(encryptor encryption.Encryption, logger *slog.Logger) (*DepositRepositoryImpl, error) {
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

	return &DepositRepositoryImpl{db: db, encryptor: encryptor, logger: logger}, nil
}

func (r *DepositRepositoryImpl) Create(ctx context.Context, deposit *domains.Deposit) error {
	fieldsToEncrypt := map[string]string{
		constants.AccountNumberField:     deposit.BankAccount.AccountNumber,
		constants.AccountHolderNameField: deposit.BankAccount.AccountHolderName,
		constants.RoutingNumberField:     deposit.BankAccount.RoutingNumber,
		constants.CardNumberField:        deposit.Card.CardNumber,
		constants.CvvField:               deposit.Card.Cvv,
	}
	encryptedFields, err := r.encryptor.EncryptMultiple(fieldsToEncrypt)
	if err != nil {
		return err
	}
	query := `
        INSERT INTO deposits (
            id, amount, currency, country, route_type, customer_id, card_number, cvv, account_number,routing_number,account_holder_name,
            status, client_reference_id, gateway_identifier, psp_response_code, psp_response_description, description,
            created_date_utc, modified_date_utc
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
    `
	_, err = r.db.ExecContext(ctx, query, deposit.Id, deposit.Amount.Value, deposit.Amount.Currency, deposit.Country,
		deposit.TransactionRouteType, deposit.Customer.Id, encryptedFields[constants.CardNumberField], encryptedFields[constants.CvvField], encryptedFields[constants.AccountNumberField],
		encryptedFields[constants.RoutingNumberField], encryptedFields[constants.AccountHolderNameField], deposit.TransactionStatus, deposit.ClientReferenceId,
		deposit.GatewayIdentifier, deposit.PspResponseCode, deposit.PspResponseDescription, deposit.Description, deposit.CreatedDateUtc, deposit.ModifiedDateUtc)
	return err
}

func (r *DepositRepositoryImpl) Update(ctx context.Context, deposit *domains.Deposit) error {
	query := `
       UPDATE deposits SET status=$2, psp_response_code=$3, psp_response_description=$4, description=$5, modified_date_utc = $6
       WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, deposit.Id, deposit.TransactionStatus, deposit.PspResponseCode, deposit.PspResponseDescription,
		deposit.Description, deposit.ModifiedDateUtc)
	return err
}

func (r *DepositRepositoryImpl) GetById(ctx context.Context, id uuid.UUID) (*domains.Deposit, error) {
	query := `
        SELECT 
            id, amount, currency, country, route_type, customer_id, card_number, cvv, account_number, routing_number, account_holder_name, 
            status, client_reference_id, gateway_identifier, psp_response_code, psp_response_description, description,
            created_date_utc, modified_date_utc
        FROM 
            deposits
        WHERE 
            id = $1
    `

	var deposit domains.Deposit
	var amountValue float64
	var customerId int64
	var currency, country, transactionStatus, routetype int
	var cardNumber, cardCvv, bankAccountNumber, bankRoutingNumber, bankAccountHolderName string

	err := r.db.QueryRow(query, id).Scan(&deposit.Id, &amountValue, &currency, &country, &routetype, &customerId,
		&cardNumber, &cardCvv, &bankAccountNumber, &bankRoutingNumber, &bankAccountHolderName,
		&transactionStatus, &deposit.ClientReferenceId, &deposit.GatewayIdentifier,
		&deposit.PspResponseCode, &deposit.PspResponseDescription, &deposit.Description,
		&deposit.CreatedDateUtc, &deposit.ModifiedDateUtc)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No deposit found with the given ID
		}
		return nil, err
	}
	parsedCurrency, err := enums.ParseCurrency(currency)
	deposit.Amount = value_objects.Amount{
		Value:    amountValue,
		Currency: parsedCurrency,
	}
	deposit.Country, _ = enums.ParseCountry(country)
	deposit.TransactionStatus, _ = enums.ParseTransactionStatus(transactionStatus)
	deposit.Customer = domains.Customer{Id: customerId}

	decryptedCardNumber, _ := r.encryptor.Decrypt(cardNumber)
	decryptedCvv, _ := r.encryptor.Decrypt(cardCvv)
	decryptedBankAccountNumber, _ := r.encryptor.Decrypt(bankAccountNumber)
	decryptedBankAccountHolderName, _ := r.encryptor.Decrypt(bankAccountHolderName)
	decryptedBankRoutingNumber, _ := r.encryptor.Decrypt(bankRoutingNumber)

	deposit.Card = value_objects.Card{
		CardNumber: decryptedCardNumber,
		Cvv:        decryptedCvv,
	}
	deposit.BankAccount = value_objects.BankAccount{
		AccountNumber:     decryptedBankAccountNumber,
		RoutingNumber:     decryptedBankRoutingNumber,
		AccountHolderName: decryptedBankAccountHolderName,
	}
	deposit.TransactionRouteType, _ = enums.ParseTransactionRouteType(routetype)
	return &deposit, nil
}
