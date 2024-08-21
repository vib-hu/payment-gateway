package controllers

import (
	callbackdtos "PaymentGateway/internal/payment-gateway-api/controllers/models/callback-dtos"
	"PaymentGateway/internal/payment-gateway-api/controllers/models/deposit-dtos"
	withdrawdtos "PaymentGateway/internal/payment-gateway-api/controllers/models/withdraw-dtos"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"PaymentGateway/internal/payment-gateway-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
)

type PaymentController struct {
	paymentService services.PaymentService
	logger         *slog.Logger
}

func NewPaymentController(service services.PaymentService, logger *slog.Logger) *PaymentController {
	return &PaymentController{paymentService: service, logger: logger}
}

// Deposit godoc
// @Summary Create a new deposit to bank or card
// @Description Creates a new deposit request and processes it.
// @Tags Deposit
// @Accept json
// @Produce json
// @Param deposit body depositdtos.DepositRequest true "Deposit request payload"
// @Success 200 {object} depositdtos.DepositResponse
// @Failure 400 {object} depositdtos.DepositResponse
// @Failure 500 {object} depositdtos.DepositResponse
// @Router /v1/payments/deposit [post]
func (pc *PaymentController) Deposit(ctx *gin.Context) {
	var req depositdtos.DepositRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorMsg := "invalid request payload"
		pc.logger.Error(errorMsg)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	pc.logger.Info("Received deposit request", "req", req)
	log.Print("outside log" + req.ClientReferenceId)
	response, err := pc.paymentService.Deposit(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	httpStatusCode := pc.getHttpStatusCodeByErrorType(response.ErrorResponse.Type)
	ctx.JSON(httpStatusCode, response)
	return
}

// Withdraw godoc
// @Summary Withdraw money using Card, ApplePay or GooglePay
// @Description Withdraw money using Card, ApplePay or GooglePay
// @Tags Withdraw
// @Accept json
// @Produce json
// @Param withdraw body withdrawdtos.WithdrawRequest true "Withdraw request payload"
// @Success 200 {object} withdrawdtos.WithdrawRequest
// @Failure 400 {object} withdrawdtos.WithdrawRequest
// @Failure 500 {object} withdrawdtos.WithdrawRequest
// @Router /v1/payments/withdraw [post]
func (pc *PaymentController) Withdraw(ctx *gin.Context) {
	var req withdrawdtos.WithdrawRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	pc.logger.Info("Received withdraw request", "req", req)
	log.Print("outside log" + req.ClientReferenceId)

	response, err := pc.paymentService.Withdraw(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	httpStatusCode := pc.getHttpStatusCodeByErrorType(response.ErrorResponse.Type)
	ctx.JSON(httpStatusCode, response)
	return
}

// Callback godoc
// @Summary Updates a deposit or withdrawal transaction
// @Description Callback endpoint for updating the deposit and withdraw transactions
// @Tags Callback
// @Accept json
// @Produce json
// @Param paymentType path string true "Payment type (deposit or withdraw)"
// @Param request body callbackdtos.CallbackRequest true "Update request payload"
// @Success 200 {object} callbackdtos.CallbackResponse
// @Failure 400 {object} callbackdtos.CallbackResponse
// @Failure 500 {object} callbackdtos.CallbackResponse
// @Router /v1/payments/{paymentType}/update [put]
func (pc *PaymentController) Callback(ctx *gin.Context) {
	paymentType := ctx.Param("paymentType")
	var req callbackdtos.CallbackRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id or status in payload."})
		return
	}

	pc.logger.Info("Received callback request", "req", req)
	paymentTypeEnum, _ := enums.ParseToPaymentType(paymentType) // ignoring error since value is already validated
	response, err := pc.paymentService.Callback(paymentTypeEnum, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	httpStatusCode := pc.getHttpStatusCodeByErrorType(response.ErrorResponse.Type)
	ctx.JSON(httpStatusCode, response)
	return
}

func (pc *PaymentController) getHttpStatusCodeByErrorType(errorType enums.ErrorType) int {
	switch errorType {
	case enums.ValidationError:
		return http.StatusBadRequest
	case enums.InternalServerError:
		return http.StatusInternalServerError
	case enums.IntegrationError:
		return http.StatusFailedDependency
	default:
		return http.StatusOK
	}
}

func ValidatePaymentType(c *gin.Context) {
	paymentType := c.Param("paymentType")
	switch paymentType {
	case string(enums.Deposit), string(enums.Withdraw):
		c.Next()
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "payment type must be deposit or withdraw."})
	}
}
