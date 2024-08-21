package gatewaya

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	gatewaymodels "PaymentGateway/internal/payment-gateway-api/integrations/gateway-a/models"
	commonmodels "PaymentGateway/internal/payment-gateway-api/integrations/models"
	"PaymentGateway/pkg/formatters"
	"PaymentGateway/pkg/protocols/http"
	"PaymentGateway/pkg/protocols/models"
	"encoding/json"
	"log/slog"
	"os"
)

type GatewayA struct {
	formatter formatters.Formatter // Json Formatter
	client    httpclient.HTTPClient
	logger    *slog.Logger
}

func NewGatewayA(formatter formatters.Formatter, protocol httpclient.HTTPClient, logger *slog.Logger) *GatewayA {
	return &GatewayA{
		formatter: formatter,
		client:    protocol,
		logger:    logger,
	}
}

func (g *GatewayA) Deposit(deposit *domains.Deposit) (commonmodels.DepositResponse, error) {
	request, err := gatewaymodels.ConvertToGatewayADepositRequest(deposit)
	if err != nil {
		g.logger.Error("error while preparing the gateway A deposit request", "deposit", deposit, err)
		return commonmodels.DepositResponse{}, err
	}

	formattedRequest, err := g.formatter.Marshal(request)
	if err != nil {
		g.logger.Error("error while serializing the gateway A deposit request", "deposit", request, err)
		return commonmodels.DepositResponse{}, err
	}

	response, err := g.client.Post("http://mock_payment_gateways:8080/api/v1/gateway-a/deposit", g.buildHeaders(),
		formattedRequest, models.DefaultResiliencyParameters("Gateway_A_Deposit"))
	if err != nil {
		g.logger.Error("error while calling the gateway A deposit endpoint", err)
		return commonmodels.DepositResponse{}, err
	}

	var gatewayResponse gatewaymodels.GatewayADepositResponse
	err = json.Unmarshal(response, &gatewayResponse)
	if err != nil {
		g.logger.Error("error while deserializing the gateway A deposit response", err)
		return commonmodels.DepositResponse{}, err
	}
	return commonmodels.DepositResponse{
		IsSuccessful:        gatewayResponse.IsSuccessful,
		ResponseCode:        gatewayResponse.ResponseCode,
		ResponseDescription: gatewayResponse.ResponseDescription,
	}, nil
}

func (g *GatewayA) Withdraw(withdraw *domains.Withdraw) (commonmodels.WithdrawResponse, error) {
	request, err := gatewaymodels.ConvertToGatewayAWithdrawRequest(withdraw)
	if err != nil {
		g.logger.Error("error while preparing the gateway A withdraw request", "request", request, err)
		return commonmodels.WithdrawResponse{}, err
	}

	formattedRequest, err := g.formatter.Marshal(request)
	if err != nil {
		g.logger.Error("error while serializing the gateway A withdraw request", "request", request, err)
		return commonmodels.WithdrawResponse{}, err
	}

	response, err := g.client.Post("http://mock_payment_gateways:8080/api/v1/gateway-a/withdraw", g.buildHeaders(),
		formattedRequest, models.DefaultResiliencyParameters("Gateway_A_Withdraw"))
	if err != nil {
		g.logger.Error("error while calling the gateway A withdraw endpoint", "error", err)
		return commonmodels.WithdrawResponse{}, err
	}

	var gatewayResponse gatewaymodels.GatewayAWithdrawResponse
	err = json.Unmarshal(response, &gatewayResponse)
	if err != nil {
		g.logger.Error("error while deserializing the gateway A withdraw response", "gatewayresponse", gatewayResponse, err)
		return commonmodels.WithdrawResponse{}, err
	}
	return commonmodels.WithdrawResponse{
		IsSuccessful:        gatewayResponse.IsSuccessful,
		ResponseCode:        gatewayResponse.ResponseCode,
		ResponseDescription: gatewayResponse.ResponseDescription,
	}, nil
}

func (g *GatewayA) GetGatewayIdentifier() string {
	return "GatewayA"
}

func (g *GatewayA) buildHeaders() map[string]string {
	headers := make(map[string]string)
	headers["API_KEY"] = os.Getenv("GATEWAY_A_API_KEY")
	headers["Content-Type"] = os.Getenv("application/json")
	return headers
}
