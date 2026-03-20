package typedsocket_test

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/mister-webhooks/sts-assume-role-proxy/typedsocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PipeConn struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func NewPipeConn() *PipeConn {
	r, w := io.Pipe()

	return &PipeConn{
		reader: r,
		writer: w,
	}
}

func (pc *PipeConn) Read(b []byte) (int, error) {
	return pc.reader.Read(b)
}

func (pc *PipeConn) Write(b []byte) (int, error) {
	i, e := pc.writer.Write(b)
	return i, e
}

func (pc *PipeConn) Close() error {
	err1 := pc.reader.Close()
	err2 := pc.writer.Close()

	return multierror.Append(err1, err2)
}

func (pc *PipeConn) LocalAddr() net.Addr {
	return nil
}

func (pc *PipeConn) RemoteAddr() net.Addr {
	return nil
}

func (pc *PipeConn) SetDeadline(time.Time) error {
	return nil
}

func (pc *PipeConn) SetReadDeadline(time.Time) error {
	return nil
}

func (pc *PipeConn) SetWriteDeadline(time.Time) error {
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
		pc := NewPipeConn()

		rtc := typedsocket.NewTypedConnection[bmstring, *bmstring](pc)
		wtc := typedsocket.NewTypedConnection[bmstring, *bmstring](pc)

		var err error

		go func() {
			err := wtc.Send(BMString("hi there"))
			require.NoError(t, err)
		}()

		x := new(bmstring)
		err = rtc.Recv(x)

		require.NoError(t, err)

		assert.Equal(t, BMString("hi there"), *x)
	})
}
