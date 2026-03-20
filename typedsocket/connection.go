package typedsocket

import (
	"encoding/json"
	"net"
)

type TypedConnection[S any, R any] struct {
	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewTypedConnection[S any, R any](conn net.Conn) *TypedConnection[S, R] {
	return &TypedConnection[S, R]{
		conn:    conn,
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
	}
}

func (tc TypedConnection[S, R]) Recv(r R) error {
	return tc.decoder.Decode(r)
}

func (tc TypedConnection[S, R]) Send(msg S) error {
	return tc.encoder.Encode(msg)
}

func (tc TypedConnection[S, R]) Close() error {
	return tc.conn.Close()
}

func (tc TypedConnection[S, R]) LocalAddr() net.Addr {
	return tc.conn.LocalAddr()
}

func (tc TypedConnection[S, R]) RemoteAddr() net.Addr {
	return tc.conn.RemoteAddr()
}

func (tc TypedConnection[S, R]) Conn() net.Conn {
	return tc.conn
}
