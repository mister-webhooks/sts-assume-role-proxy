package typedsocket_test

import (
	"encoding/json"
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

type JString struct {
	String *string
}

func NewJString(str string) JString {
	return JString{String: &str}
}

func (s JString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String)
}

func (s *JString) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.String)
}

func TestTypedConnection(t *testing.T) {
	t.Run("happy path send / recv works", func(t *testing.T) {
		pc := NewPipeConn()

		rtc := typedsocket.NewTypedConnection[JString, *JString](pc)
		wtc := typedsocket.NewTypedConnection[JString, *JString](pc)

		var err error

		go func() {
			err := wtc.Send(NewJString("hi there"))
			require.NoError(t, err)
		}()

		x := new(JString)
		err = rtc.Recv(x)

		require.NoError(t, err)

		assert.Equal(t, NewJString("hi there"), *x)
	})
}
