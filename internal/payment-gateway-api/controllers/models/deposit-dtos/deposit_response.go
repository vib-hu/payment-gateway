package depositdtos

import (
	"PaymentGateway/internal/payment-gateway-api/controllers/models/common-dtos"
	"github.com/google/uuid"
)

type DepositResponse struct {
	IsSuccessful  bool                     `json:"is_successful"`
	TransactionId uuid.UUID                `json:"transaction_id"`
	ErrorResponse commondtos.ErrorResponse `json:"error_response"`
}

func NewDepositResponse(transactionId uuid.UUID) *DepositResponse {
	return &DepositResponse{
		TransactionId: transactionId,
	}
}

func (d *DepositResponse) Success() *DepositResponse {
	d.IsSuccessful = true
	return d
}

func (d *DepositResponse) Failed(errorResponse commondtos.ErrorResponse) *DepositResponse {
	d.IsSuccessful = false
	d.ErrorResponse = errorResponse
	return d
}
