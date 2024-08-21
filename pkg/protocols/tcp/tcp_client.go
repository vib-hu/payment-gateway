package tcpclient

import (
	"context"
)

type HTTPClient interface {
	Send(ctx context.Context, address string, request []byte) ([]byte, error)
}
