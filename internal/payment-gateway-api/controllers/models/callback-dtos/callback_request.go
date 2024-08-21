package callbackdtos

import (
	"PaymentGateway/internal/payment-gateway-api/enums"
	"github.com/google/uuid"
)

type CallbackRequest struct {
	TransactionId     uuid.UUID               `json:"transaction_id"`
	TransactionStatus enums.TransactionStatus `json:"transaction_status"`
}
