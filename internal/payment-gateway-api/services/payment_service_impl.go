package services

import (
	callbackdtos "PaymentGateway/internal/payment-gateway-api/controllers/models/callback-dtos"
	commondtos "PaymentGateway/internal/payment-gateway-api/controllers/models/common-dtos"
	"PaymentGateway/internal/payment-gateway-api/controllers/models/deposit-dtos"
	withdrawdtos "PaymentGateway/internal/payment-gateway-api/controllers/models/withdraw-dtos"
	"PaymentGateway/internal/payment-gateway-api/domains"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"PaymentGateway/internal/payment-gateway-api/integrations"
	"PaymentGateway/internal/payment-gateway-api/repositories"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

type PaymentServiceImpl struct {
	gateways           map[enums.Country]integrations.Gateway
	depositRepository  repositories.DepositRepository
	withdrawRepository repositories.WithdrawRepository
	logger             *slog.Logger
}

func NewPaymentService(gateways map[enums.Country]integrations.Gateway,
	depositRepository repositories.DepositRepository, withdrawRepository repositories.WithdrawRepository, logger *slog.Logger) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		gateways:           gateways,
		depositRepository:  depositRepository,
		withdrawRepository: withdrawRepository,
		logger:             logger,
	}
}

func (ps *PaymentServiceImpl) Deposit(request depositdtos.DepositRequest) (*depositdtos.DepositResponse, error) {
	transactionId := uuid.New()
	response := depositdtos.NewDepositResponse(transactionId)

	// assuming gateways are configured per country
	gateway, exists := ps.gateways[request.Country]

	if !exists {
		errorMsg := fmt.Sprintf("no gateway configured in the requested country: %d", request.Country)
		ps.logger.Error(errorMsg, "request", request)
		errorResponse := ps.buildError(enums.ValidationError, errorMsg)
		return response.Failed(errorResponse), nil
	}

	deposit, err := domains.NewDeposit(transactionId, request.Amount.ToDomain(), request.Country, request.TransactionRouteType,
		request.Customer.ToDomain(), request.CardDetails.ToDomain(), request.BankDetails.ToDomain(),
		gateway.GetGatewayIdentifier(), request.ClientReferenceId)
	if err != nil {
		errorMsg := fmt.Sprintf("validation failed: %s", err.Error())
		ps.logger.Warn(errorMsg, "request", request)
		errorResponse := ps.buildError(enums.ValidationError, errorMsg)
		return response.Failed(errorResponse), nil
	}

	err = ps.depositRepository.Create(context.Background(), deposit)
	if err != nil {
		errorMsg := fmt.Sprintf("error while creating the deposit transaction in db: %s", err.Error())
		ps.logger.Error(errorMsg, "request", request)

		// hiding actual details to be exposed to the client due to security issue
		errorResponse := ps.buildError(enums.InternalServerError, "internal server error")
		return response.Failed(errorResponse), nil
	}

	depositResponse, err := gateway.Deposit(deposit)
	if err != nil {
		errorMsg := fmt.Sprintf("error while invoking gateway %s deposit endpoint: %s", gateway.GetGatewayIdentifier(), err.Error())
		ps.logger.Error(errorMsg, "request", request)

		// hiding actual details to be exposed to the client due to security issue
		errorResponse := ps.buildError(enums.IntegrationError, "payment integration error")
		deposit.FailedDueToInternalError(err.Error())
		ps.depositRepository.Update(context.Background(), deposit)
		return response.Failed(errorResponse), nil
	}

	if depositResponse.IsSuccessful {
		deposit.Success(depositResponse.ResponseCode, depositResponse.ResponseDescription)
		ps.depositRepository.Update(context.Background(), deposit)
		return response.Success(), nil
	}

	deposit.FailedDueToPayment(depositResponse.ResponseCode, depositResponse.ResponseDescription)
	ps.depositRepository.Update(context.Background(), deposit)
	errorResponse := ps.buildError(enums.PaymentError, depositResponse.ResponseDescription)
	return response.Failed(errorResponse), nil
}

func (ps *PaymentServiceImpl) Withdraw(request withdrawdtos.WithdrawRequest) (*withdrawdtos.WithdrawResponse, error) {
	transactionId := uuid.New()
	response := withdrawdtos.NewWithdrawResponse(transactionId)

	// assuming gateways are configured per country
	gateway, exists := ps.gateways[request.Country]

	if !exists {
		errorMsg := fmt.Sprintf("no gateway configured in the requested country: %d", request.Country)
		ps.logger.Error(errorMsg, "request", request)
		errorResponse := ps.buildError(enums.ValidationError, "No Gateway configured in the requested country")
		return response.Failed(errorResponse), nil
	}

	withdraw, err := domains.NewWithdraw(transactionId, request.Amount.ToDomain(), request.Country,
		request.Customer.ToDomain(), request.PaymentMethod, request.CardDetails.ToDomain(), request.ApplePayDetails.ToDomain(),
		request.GooglePayDetails.ToDomain(), gateway.GetGatewayIdentifier(), request.ClientReferenceId, request.BillingAddress)

	if err != nil {
		errorMsg := fmt.Sprintf("validation failed: %s", err.Error())
		ps.logger.Warn(errorMsg, "request", request)
		errorResponse := ps.buildError(enums.ValidationError, errorMsg)
		return response.Failed(errorResponse), nil
	}

	err = ps.withdrawRepository.Create(context.Background(), withdraw)
	if err != nil {
		errorMsg := fmt.Sprintf("error while creating the withdraw transaction in db: %s", err.Error())
		ps.logger.Error(errorMsg, "request", request)

		// hiding actual details to be exposed to the client due to security issue
		errorResponse := ps.buildError(enums.InternalServerError, "internal server error")
		return response.Failed(errorResponse), nil
	}

	withdrawResponse, err := gateway.Withdraw(withdraw)
	if err != nil {
		errorMsg := fmt.Sprintf("error while invoking gateway %s withdraw endpoint: %s", gateway.GetGatewayIdentifier(), err.Error())
		ps.logger.Error(errorMsg, "request", request)

		// hiding actual details to be exposed to the client due to security issue
		errorResponse := ps.buildError(enums.IntegrationError, "payment integration error")
		withdraw.FailedDueToInternalError(err.Error())
		ps.withdrawRepository.Update(context.Background(), withdraw)
		return response.Failed(errorResponse), nil
	}

	if withdrawResponse.IsSuccessful {
		withdraw.Success(withdrawResponse.ResponseCode, withdrawResponse.ResponseDescription)
		ps.withdrawRepository.Update(context.Background(), withdraw)
		return response.Success(), nil
	}

	withdraw.FailedDueToPayment(withdrawResponse.ResponseCode, withdrawResponse.ResponseDescription)
	ps.withdrawRepository.Update(context.Background(), withdraw)
	errorResponse := ps.buildError(enums.PaymentError, withdrawResponse.ResponseDescription)
	return response.Failed(errorResponse), nil
}

func (ps *PaymentServiceImpl) Callback(paymentType enums.PaymentType, request callbackdtos.CallbackRequest) (*callbackdtos.CallbackResponse, error) {
	response := callbackdtos.NewCallbackResponse()
	switch paymentType {
	case enums.Deposit:
		return ps.updateDepositTransaction(request, response)
	case enums.Withdraw:
		return ps.updateWithdrawTransaction(request, response)
	default:
		return response.Failed(ps.buildError(enums.ValidationError, "invalid payment type")), nil
	}
}

func (ps *PaymentServiceImpl) updateDepositTransaction(request callbackdtos.CallbackRequest, response *callbackdtos.CallbackResponse) (*callbackdtos.CallbackResponse, error) {
	deposit, err := ps.depositRepository.GetById(context.Background(), request.TransactionId)
	if err != nil {
		errorMsg := fmt.Sprintf("error while fetching deposit transaction from db: %s", err.Error())
		ps.logger.Error(errorMsg, "request", request)

		// hiding actual details to be exposed to the client due to security issue
		errorResponse := ps.buildError(enums.InternalServerError, "internal server error")
		return response.Failed(errorResponse), nil
	}
	if deposit == nil {
		errorMsg := fmt.Sprintf("no deposit transaction found with given id: %s", request.TransactionId)
		ps.logger.Warn(errorMsg, "request", request)
		errorResponse := ps.buildError(enums.ValidationError, errorMsg)
		return response.Failed(errorResponse), nil
	}

	deposit.UpdateStatus(request.TransactionStatus)
	err = ps.depositRepository.Update(context.Background(), deposit)
	if err != nil {
		errorMsg := fmt.Sprintf("error while updating deposit transaction in db: %s", err.Error())
		ps.logger.Error(errorMsg, "request", request)

		// hiding actual details to be exposed to the client due to security issue
		errorResponse := ps.buildError(enums.InternalServerError, "internal server error")
		return response.Failed(errorResponse), nil
	}
	return response.Success(), nil
}

func (ps *PaymentServiceImpl) updateWithdrawTransaction(request callbackdtos.CallbackRequest, response *callbackdtos.CallbackResponse) (*callbackdtos.CallbackResponse, error) {
	withdraw, err := ps.withdrawRepository.GetById(context.Background(), request.TransactionId)
	if err != nil {
		errorMsg := fmt.Sprintf("error while fetching withdraw transaction from db: %s", err.Error())
		ps.logger.Error(errorMsg, "request", request)

		// hiding actual details to be exposed to the client due to security issue
		errorResponse := ps.buildError(enums.InternalServerError, "internal server error")
		return response.Failed(errorResponse), nil
	}
	if withdraw == nil {
		errorMsg := fmt.Sprintf("no withdraw transaction found with given id: %s", request.TransactionId)
		ps.logger.Warn(errorMsg, "request", request)
		errorResponse := ps.buildError(enums.ValidationError, errorMsg)
		return response.Failed(errorResponse), nil
	}

	withdraw.UpdateStatus(request.TransactionStatus)
	err = ps.withdrawRepository.Update(context.Background(), withdraw)
	if err != nil {
		errorMsg := fmt.Sprintf("error while updating withdraw transaction in db: %s", err.Error())
		ps.logger.Error(errorMsg, "request", request)

		// hiding actual details to be exposed to the client due to security issue
		errorResponse := ps.buildError(enums.InternalServerError, "internal server error")
		return response.Failed(errorResponse), nil
	}
	return response.Success(), nil
}

func (ps *PaymentServiceImpl) buildError(errorType enums.ErrorType, errorMessage string) commondtos.ErrorResponse {
	return commondtos.ErrorResponse{
		Type:    errorType,
		Message: errorMessage,
	}
}
