package callbackdtos

import (
	commondtos "PaymentGateway/internal/payment-gateway-api/controllers/models/common-dtos"
)

type CallbackResponse struct {
	IsSuccessful  bool                     `json:"is_successful"`
	ErrorResponse commondtos.ErrorResponse `json:"error_response"`
}

func NewCallbackResponse() *CallbackResponse {
	return &CallbackResponse{}
}

func (c *CallbackResponse) Success() *CallbackResponse {
	c.IsSuccessful = true
	return c
}

func (c *CallbackResponse) Failed(errorResponse commondtos.ErrorResponse) *CallbackResponse {
	c.IsSuccessful = false
	c.ErrorResponse = errorResponse
	return c
}
