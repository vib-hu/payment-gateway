package services

import (
	commondtos "PaymentGateway/internal/payment-gateway-api/controllers/models/common-dtos"
	"PaymentGateway/internal/payment-gateway-api/controllers/models/deposit-dtos"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"PaymentGateway/internal/payment-gateway-api/integrations"
	"PaymentGateway/internal/payment-gateway-api/integrations/models"
	integrationdtos "PaymentGateway/internal/payment-gateway-api/integrations/models"
	"PaymentGateway/internal/payment-gateway-api/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

func TestPaymentServiceImpl_Deposit_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockGateway := mocks.NewMockGateway(ctrl)
	mockDepositRepo := mocks.NewMockDepositRepository(ctrl)
	mockWithdrawRepo := mocks.NewMockWithdrawRepository(ctrl)

	// Create a logger
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	// Initialize the service with mocks
	ps := NewPaymentService(
		map[enums.Country]integrations.Gateway{
			enums.CountryUSA: mockGateway,
		},
		mockDepositRepo,
		mockWithdrawRepo,
		logger,
	)

	// Test data
	request := depositdtos.DepositRequest{
		Country:              enums.CountryUSA,
		Amount:               commondtos.AmountRequest{Value: 100.0, Currency: enums.CurrencyUSA},
		TransactionRouteType: enums.RouteToBank,
		BankDetails: commondtos.BankDetailsRequest{
			AccountNumber:     "123",
			AccountHolderName: "vibhu mishra",
			RoutingNumber:     "route-123",
		},
		Customer: commondtos.CustomerRequest{
			Id: 12,
		},
		ClientReferenceId: "client-reference-123",
	}

	mockDepositRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	mockGateway.EXPECT().Deposit(gomock.Any()).Return(models.DepositResponse{IsSuccessful: true}, nil)
	mockGateway.EXPECT().GetGatewayIdentifier().Return("test")
	mockDepositRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	_, err := ps.Deposit(request)

	assert.NoError(t, err)
}

func TestPaymentServiceImpl_Deposit_No_Gateway(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockGateway := mocks.NewMockGateway(ctrl)
	mockDepositRepo := mocks.NewMockDepositRepository(ctrl)
	mockWithdrawRepo := mocks.NewMockWithdrawRepository(ctrl)

	// Create a logger
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	// Initialize the service with mocks
	ps := NewPaymentService(
		map[enums.Country]integrations.Gateway{
			enums.CountryUSA: mockGateway,
		},
		mockDepositRepo,
		mockWithdrawRepo,
		logger,
	)

	// Test data
	request := depositdtos.DepositRequest{
		Country:              enums.CountryUSA,
		Amount:               commondtos.AmountRequest{Value: 100.0, Currency: enums.CurrencyUSA},
		TransactionRouteType: enums.RouteToBank,
		BankDetails: commondtos.BankDetailsRequest{
			AccountNumber:     "123",
			AccountHolderName: "vibhu mishra",
			RoutingNumber:     "route-123",
		},
		Customer: commondtos.CustomerRequest{
			Id: 12,
		},
		ClientReferenceId: "client-reference-123",
	}

	request.Country = enums.CountryUK // No gateway configured for UK
	resp, err := ps.Deposit(request)
	assert.NoError(t, err)
	assert.False(t, resp.IsSuccessful)
}

func TestPaymentServiceImpl_Deposit_Gateway_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockGateway := mocks.NewMockGateway(ctrl)
	mockDepositRepo := mocks.NewMockDepositRepository(ctrl)
	mockWithdrawRepo := mocks.NewMockWithdrawRepository(ctrl)

	// Create a logger
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	// Initialize the service with mocks
	ps := NewPaymentService(
		map[enums.Country]integrations.Gateway{
			enums.CountryUSA: mockGateway,
		},
		mockDepositRepo,
		mockWithdrawRepo,
		logger,
	)

	// Test data
	request := depositdtos.DepositRequest{
		Country:              enums.CountryUSA,
		Amount:               commondtos.AmountRequest{Value: 100.0, Currency: enums.CurrencyUSA},
		TransactionRouteType: enums.RouteToBank,
		BankDetails: commondtos.BankDetailsRequest{
			AccountNumber:     "123",
			AccountHolderName: "vibhu mishra",
			RoutingNumber:     "route-123",
		},
		Customer: commondtos.CustomerRequest{
			Id: 12,
		},
		ClientReferenceId: "client-reference-123",
	}

	mockDepositRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	mockGateway.EXPECT().Deposit(gomock.Any()).Return(integrationdtos.DepositResponse{}, errors.New("gateway error"))
	mockGateway.EXPECT().GetGatewayIdentifier().Return("test").AnyTimes()
	mockDepositRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := ps.Deposit(request)
	assert.NoError(t, err)
	assert.False(t, resp.IsSuccessful)
}
