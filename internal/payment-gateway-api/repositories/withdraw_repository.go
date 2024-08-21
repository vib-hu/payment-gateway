package repositories

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=../mocks/mock_withdraw_repository.go -package=mocks PaymentGateway/internal/payment-gateway-api/repositories WithdrawRepository
type WithdrawRepository interface {
	Create(ctx context.Context, withdraw *domains.Withdraw) error
	Update(ctx context.Context, withdraw *domains.Withdraw) error
	GetById(ctx context.Context, id uuid.UUID) (*domains.Withdraw, error)
}
