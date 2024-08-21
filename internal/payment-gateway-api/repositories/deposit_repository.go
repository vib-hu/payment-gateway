package repositories

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=../mocks/mock_deposit_repository.go -package=mocks PaymentGateway/internal/payment-gateway-api/repositories DepositRepository
type DepositRepository interface {
	Create(ctx context.Context, deposit *domains.Deposit) error
	Update(ctx context.Context, deposit *domains.Deposit) error
	GetById(ctx context.Context, id uuid.UUID) (*domains.Deposit, error)
}
