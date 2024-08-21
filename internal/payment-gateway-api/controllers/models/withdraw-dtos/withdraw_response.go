package withdrawdtos

import (
	commondtos "PaymentGateway/internal/payment-gateway-api/controllers/models/common-dtos"
	"github.com/google/uuid"
)

type WithdrawResponse struct {
	IsSuccessful  bool                     `json:"is_successful"`
	TransactionId uuid.UUID                `json:"transaction_id"`
	ErrorResponse commondtos.ErrorResponse `json:"error_response"`
}

func NewWithdrawResponse(transactionId uuid.UUID) *WithdrawResponse {
	return &WithdrawResponse{
		TransactionId: transactionId,
	}
}

func (d *WithdrawResponse) Success() *WithdrawResponse {
	d.IsSuccessful = true
	return d
}

func (d *WithdrawResponse) Failed(errorResponse commondtos.ErrorResponse) *WithdrawResponse {
	d.IsSuccessful = false
	d.ErrorResponse = errorResponse
	return d
}
