package gatewayb

import (
	"PaymentGateway/internal/payment-gateway-api/domains"
	gatewaymodels "PaymentGateway/internal/payment-gateway-api/integrations/gateway-b/models"
	commonmodels "PaymentGateway/internal/payment-gateway-api/integrations/models"
	"PaymentGateway/pkg/formatters"
	"PaymentGateway/pkg/protocols/http"
	"PaymentGateway/pkg/protocols/models"
	"log/slog"
	"os"
)

type GatewayB struct {
	formatter formatters.Formatter // SOAP-XML Formatter
	client    httpclient.HTTPClient
	logger    *slog.Logger
}

func NewGatewayB(formatter formatters.Formatter, protocol httpclient.HTTPClient, logger *slog.Logger) *GatewayB {
	return &GatewayB{
		formatter: formatter,
		client:    protocol,
		logger:    logger,
	}
}

func (g *GatewayB) Deposit(deposit *domains.Deposit) (commonmodels.DepositResponse, error) {
	request, err := gatewaymodels.ConvertToGatewayBDepositRequest(deposit)
	if err != nil {
		g.logger.Error("error while preparing the gateway B deposit request", "deposit", deposit, err)
		return commonmodels.DepositResponse{}, err
	}

	formattedRequest, err := g.formatter.Marshal(request)
	if err != nil {
		g.logger.Error("error while serializing the gateway B deposit request", "request", request, err)
		return commonmodels.DepositResponse{}, err
	}

	rawResponse, err := g.client.Post("http://mock_payment_gateways:8080/api/v1/gateway-b/deposit", g.buildHeaders(),
		formattedRequest, models.DefaultResiliencyParameters("Gateway_B_Deposit"))
	if err != nil {
		g.logger.Error("error while calling the gateway B deposit endpoint", "error", err.Error())
		return commonmodels.DepositResponse{}, err
	}
	var envelop gatewaymodels.GatewayBDepositEnvelope
	err = g.formatter.Unmarshal(rawResponse, &envelop)
	if err != nil {
		g.logger.Error("error while deserializing the gateway B deposit response", "error", err.Error())
		return commonmodels.DepositResponse{}, err
	}

	response := envelop.Body.GatewayBDepositResponse
	return commonmodels.DepositResponse{
		IsSuccessful:        response.IsSuccessful,
		ResponseCode:        response.ResponseCode,
		ResponseDescription: response.ResponseDescription,
	}, nil
}

func (g *GatewayB) Withdraw(withdraw *domains.Withdraw) (commonmodels.WithdrawResponse, error) {
	request, err := gatewaymodels.ConvertToGatewayBWithdrawRequest(withdraw)
	if err != nil {
		g.logger.Error("error while preparing the gateway B withdraw request", "withdraw", withdraw, err)
		return commonmodels.WithdrawResponse{}, err
	}
	formattedRequest, err := g.formatter.Marshal(request)
	if err != nil {
		g.logger.Error("error while serializing the gateway B withdraw request", "request", request, err)
		return commonmodels.WithdrawResponse{}, err
	}
	rawResponse, err := g.client.Post("http://mock_payment_gateways:8080/api/v1/gateway-b/withdraw", g.buildHeaders(),
		formattedRequest, models.DefaultResiliencyParameters("Gateway_B_Withdraw"))
	if err != nil {
		g.logger.Error("error while calling the gateway B withdraw endpoint", err)
		return commonmodels.WithdrawResponse{}, err
	}
	var envelop gatewaymodels.GatewayBWithdrawEnvelope
	err = g.formatter.Unmarshal(rawResponse, &envelop)
	if err != nil {
		g.logger.Error("error while deserializing the gateway B withdraw response", err)
		return commonmodels.WithdrawResponse{}, err
	}
	response := envelop.Body.GatewayBWithdrawResponse
	return commonmodels.WithdrawResponse{
		IsSuccessful:        response.IsSuccessful,
		ResponseCode:        response.ResponseCode,
		ResponseDescription: response.ResponseDescription,
	}, nil
}

func (g *GatewayB) GetGatewayIdentifier() string {
	return "GatewayB"
}

func (g *GatewayB) buildHeaders() map[string]string {
	headers := make(map[string]string)
	headers["API_KEY"] = os.Getenv("GATEWAY_B_API_KEY")
	headers["Content-Type"] = os.Getenv("application/soap+xml")
	return headers
}
