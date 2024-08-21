FROM golang:1.22.6-alpine as build

WORKDIR /app
COPY . .

# Install Swag CLI tool
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
RUN swag init  --parseDependency --parseInternal --dir cmd/payment-gateway-api/,internal/payment-gateway-api/ -o ./docs/swagger

#Build
RUN go build -o payment-gateway-api ./cmd/payment-gateway-api

FROM alpine as runtime
COPY --from=build /app/payment-gateway-api /app/cmd/payment-gateway-api/settings** ./

# adding curl for health checks
RUN apk --no-cache add curl

EXPOSE 80
CMD [ "./payment-gateway-api" ]
