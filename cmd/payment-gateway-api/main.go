package main

import (
	_ "PaymentGateway/docs/swagger"
	"PaymentGateway/internal/payment-gateway-api/config"
	"PaymentGateway/internal/payment-gateway-api/controllers"
	"PaymentGateway/internal/payment-gateway-api/enums"
	"PaymentGateway/internal/payment-gateway-api/integrations"
	gatewaya "PaymentGateway/internal/payment-gateway-api/integrations/gateway-a"
	gatewayb "PaymentGateway/internal/payment-gateway-api/integrations/gateway-b"
	"PaymentGateway/internal/payment-gateway-api/repositories"
	"PaymentGateway/internal/payment-gateway-api/services"
	"PaymentGateway/pkg/encryption"
	"PaymentGateway/pkg/formatters"
	httpclient "PaymentGateway/pkg/protocols/http"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// @title           Payment Gateway API
// @version         1.0
// @description     This is a Payment Gateway API supporting Gateway A and B

// @contact.name   Vibhu Mishra
// @contact.url    https://www.linkedin.com/in/vibhu-mishra-07532236/
// @contact.email  vibhumishra808@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
func main() {

	opts := &slog.HandlerOptions{
		Level: getLogLevel(),
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	jsonFormatter := formatters.NewJsonFormatter()
	soapFormatter := formatters.NewSoapFormatter()
	httpclient := httpclient.NewResilientHttpClient()
	gatewayA := gatewaya.NewGatewayA(jsonFormatter, httpclient, logger)
	gatewayB := gatewayb.NewGatewayB(soapFormatter, httpclient, logger)

	gatewayMap := make(map[enums.Country]integrations.Gateway)
	gatewayMap[enums.CountryUSA] = gatewayA
	gatewayMap[enums.CountryUK] = gatewayB
	encryptor := encryption.NewAesEncryption()
	depositRepository, err := repositories.NewDepositRepository(encryptor, logger)
	if err != nil {
		logger.Error("error while registering deposit repository", err)
	}

	withdrawRepository, err := repositories.NewWithdrawRepository(encryptor, logger)
	if err != nil {
		logger.Error("error while registering withdraw repository", err)
	}

	paymentService := services.NewPaymentService(gatewayMap, depositRepository, withdrawRepository, logger)
	paymentController := controllers.NewPaymentController(paymentService, logger)

	router := gin.Default()
	registerSwagger(router)
	registerHealth(router)
	registerV1Controllers(router, paymentController)
	config := initConfig()
	srv := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			err = fmt.Errorf("HTTP ListenAndServe  server : %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), config.Server.HttpServerShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		err = fmt.Errorf("error while shutdown payment-gateway-api : %w", err)
		return
	}

	fmt.Sprintf("Graceful shutdown within %v", 20)
}

func registerV1Controllers(router *gin.Engine, paymentController *controllers.PaymentController) {
	v1 := router.Group("/v1")
	{
		paymentRoutesV1 := v1.Group("/payments")
		{
			paymentRoutesV1.POST("/deposit", paymentController.Deposit)
			paymentRoutesV1.POST("/withdraw", paymentController.Withdraw)
			paymentRoutesV1.PUT("/:paymentType/update", controllers.ValidatePaymentType, paymentController.Callback)
		}
	}
}

func registerHealth(router *gin.Engine) {
	router.GET("/health", health)
}

func registerSwagger(router *gin.Engine) {
	env := os.Getenv("ENV")
	//preventing swagger expose in prod environment
	if strings.Compare(env, "qa") == 0 || strings.Compare(env, "local") == 0 {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

// health godoc
// @Summary health endpoint
// @Schemes
// @Description health endpoint
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /health [get]
func health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Health Check Ok")
}

func initConfig() *config.Config {
	env := os.Getenv("ENV")
	viper.SetConfigName(fmt.Sprintf("%s_%s", "settings", env))
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv()

	var config config.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &config
}

func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default log level
	}
}
