all: lint

.PHONY: build

lint:
		scripts/lint
lint_fix:
		scripts/lint -v --fix

start-components:
		cd docker && docker-compose -f docker-compose.base.yml -f docker-compose.payment-gateway-api.yml build payment_gateway_api && \
                     docker-compose -f docker-compose.base.yml -f docker-compose.payment-gateway-api.yml up -d

stop-components:
		cd docker && docker-compose -f docker-compose.base.yml -f docker-compose.payment-gateway-api.yml down

