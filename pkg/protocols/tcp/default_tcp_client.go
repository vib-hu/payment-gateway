package tcpclient

import (
	"context"
	"net"
	"time"
)

type DefaultTcpClient struct {
}

func NewDefaultTcpClient() *DefaultTcpClient {
	return &DefaultTcpClient{}
}

func (p *DefaultTcpClient) Send(ctx context.Context, address string, data []byte) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:n], nil
}
