package typedsocket_test

import (
	"net"
	"testing"
	"time"

	"github.com/mister-webhooks/sts-assume-role-proxy/typedsocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ChanConn struct {
	c chan byte
}

func NewChanConnPair() (*ChanConn, *ChanConn) {
	c := make(chan byte, 1024*1024)

	return &ChanConn{
			c,
		}, &ChanConn{
			c,
		}
}

func (cc *ChanConn) Read(b []byte) (int, error) {
	i := 0

	for ; i < cap(b); i++ {
		b[i] = <-cc.c
	}

	return i, nil
}

func (cc *ChanConn) Write(b []byte) (int, error) {
	i := 0

	for ; i < len(b); i++ {
		cc.c <- b[i]
	}

	return i, nil
}

func (cc *ChanConn) Close() error {
	close(cc.c)
	return nil
}

func (cc *ChanConn) LocalAddr() net.Addr {
	return nil
}

func (cc *ChanConn) RemoteAddr() net.Addr {
	return nil
}

func (cc *ChanConn) SetDeadline(time.Time) error {
	return nil
}

func (cc *ChanConn) SetReadDeadline(time.Time) error {
	return nil
}

func (cc *ChanConn) SetWriteDeadline(time.Time) error {
	return nil
}

type bmstring struct {
	*string
}

func BMString(str string) bmstring {
	return bmstring{string: &str}
}

func (s bmstring) MarshalBinary() ([]byte, error) {
	return []byte(*s.string), nil
}

func (s *bmstring) UnmarshalBinary(data []byte) error {
	strdat := string(data)
	s.string = &strdat

	return nil
}

func TestTypedConnection(t *testing.T) {
	t.Run("happy path send / recv works", func(t *testing.T) {
		r, w := NewChanConnPair()

		rtc := typedsocket.NewTypedConnection[bmstring, *bmstring](r)
		wtc := typedsocket.NewTypedConnection[bmstring, *bmstring](w)

		err := wtc.Send(BMString("hi there"))

		require.NoError(t, err)

		x := new(bmstring)
		err = rtc.Recv(x)

		require.NoError(t, err)

		assert.Equal(t, BMString("hi there"), *x)
	})
}
