package typedsocket

import (
	"encoding"
	"fmt"
	"net"
)

type TypedClient[Q encoding.BinaryMarshaler, R encoding.BinaryUnmarshaler] struct {
	*TypedConnection[Q, R]
}

func Dial[Q encoding.BinaryMarshaler, R encoding.BinaryUnmarshaler](network string, address string) (*TypedClient[Q, R], error) {
	conn, err := net.Dial(network, address)

	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}

	return &TypedClient[Q, R]{
		TypedConnection: NewTypedConnection[Q, R](conn),
	}, nil
}
