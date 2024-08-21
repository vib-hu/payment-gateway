package commondtos

import "PaymentGateway/internal/payment-gateway-api/enums"

type ErrorResponse struct {
	Type    enums.ErrorType `json:"type"`
	Message string          `json:"message"`
}
