package typedsocket

import (
	"fmt"
	"net"
)

type TypedClient[Q any, R any] struct {
	*TypedConnection[Q, R]
}

func Dial[Q any, R any](network string, address string) (*TypedClient[Q, R], error) {
	conn, err := net.Dial(network, address)

	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}

	return &TypedClient[Q, R]{
		TypedConnection: NewTypedConnection[Q, R](conn),
	}, nil
}
