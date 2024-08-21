package httpclient

import "PaymentGateway/pkg/protocols/models"

type HTTPClient interface {
	Post(endpoint string, headers map[string]string, request []byte, resiliencyParameters models.ResiliencyParameters) ([]byte, error)
}
