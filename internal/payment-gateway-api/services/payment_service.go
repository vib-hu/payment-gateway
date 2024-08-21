package services

import (
	callbackdtos "PaymentGateway/internal/payment-gateway-api/controllers/models/callback-dtos"
	"PaymentGateway/internal/payment-gateway-api/controllers/models/deposit-dtos"
	withdrawdtos "PaymentGateway/internal/payment-gateway-api/controllers/models/withdraw-dtos"
	"PaymentGateway/internal/payment-gateway-api/enums"
)

//go:generate mockgen -destination=../mocks/mock_payment_service.go -package=mocks PaymentGateway/internal/payment-gateway-api/services PaymentService
type PaymentService interface {
	Deposit(request depositdtos.DepositRequest) (*depositdtos.DepositResponse, error)
	Withdraw(request withdrawdtos.WithdrawRequest) (*withdrawdtos.WithdrawResponse, error)
	Callback(paymentType enums.PaymentType, request callbackdtos.CallbackRequest) (*callbackdtos.CallbackResponse, error)
}
