package typedsocket

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"net"
)

type TypedConnection[S encoding.BinaryMarshaler, R encoding.BinaryUnmarshaler] struct {
	conn net.Conn
}

func NewTypedConnection[S encoding.BinaryMarshaler, R encoding.BinaryUnmarshaler](conn net.Conn) *TypedConnection[S, R] {
	return &TypedConnection[S, R]{
		conn,
	}
}

func (tc TypedConnection[S, R]) Recv(r R) error {
	lenbuf := make([]byte, 4)

	n, err := tc.conn.Read(lenbuf)

	if err != nil {
		return err
	}

	if n != len(lenbuf) {
		panic("error: short read of message length!")
	}

	len := binary.BigEndian.Uint32(lenbuf)

	databuf := make([]byte, len)

	n, err = tc.conn.Read(databuf)

	if err != nil {
		return err
	}

	if n != int(len) {
		panic("error: short read of databuf!")
	}

	err = r.UnmarshalBinary(databuf)

	if err != nil {
		return fmt.Errorf("could not decode message: %w", err)
	}

	return nil
}

func (tc TypedConnection[S, R]) Send(msg S) error {
	enc, err := msg.MarshalBinary()

	if err != nil {
		return err
	}

	buf := make([]byte, 0, 4+len(enc))
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(enc)))
	buf = append(buf, enc...)

	n, err := tc.conn.Write(buf)

	if err != nil {
		return err
	}

	if n != 4+len(enc) {
		panic(fmt.Sprintf("error: short write of encoded message: %d expected, %d actual", 4+len(enc), n))
	}

	return nil
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
