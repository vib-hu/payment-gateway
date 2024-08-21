package integrations

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	"PaymentGateway/internal/payment-gateway-api/integrations/models"
)

//go:generate mockgen -destination=../mocks/mock_gateway.go -package=mocks PaymentGateway/internal/payment-gateway-api/integrations Gateway
type Gateway interface {
	GetGatewayIdentifier() string
	Deposit(deposit *domains.Deposit) (models.DepositResponse, error)
	Withdraw(withdraw *domains.Withdraw) (models.WithdrawResponse, error)
}
