package protocol_test

import (
	"testing"

	"github.com/mister-webhooks/sts-assume-role-proxy/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLVString(t *testing.T) {
	t.Run("encoding round-trips for an empty string", func(t *testing.T) {
		lvs := protocol.LVString("")
		encoded, err := lvs.MarshalBinary()
		decoded := new(protocol.LVString)

		require.NoError(t, err)

		err = decoded.UnmarshalBinary(encoded)

		assert.NoError(t, err)
		assert.Equal(t, lvs, *decoded)
	})

	t.Run("encoding length is correct for an empty string", func(t *testing.T) {
		lvs := protocol.LVString("")
		encoded, err := lvs.MarshalBinary()

		require.NoError(t, err)

		assert.Equal(t, len(encoded), lvs.EncodedLength())
	})

	t.Run("encoding round-trips for a string", func(t *testing.T) {
		lvs := protocol.LVString("hello world")
		encoded, err := lvs.MarshalBinary()
		decoded := new(protocol.LVString)

		require.NoError(t, err)

		err = decoded.UnmarshalBinary(encoded)

		assert.NoError(t, err)
		assert.Equal(t, lvs, *decoded)
	})

	t.Run("encoding length is correct for a non-empty string", func(t *testing.T) {
		lvs := protocol.LVString("i am the very model of a modern major general")
		encoded, err := lvs.MarshalBinary()

		require.NoError(t, err)

		assert.Equal(t, len(encoded), lvs.EncodedLength())
	})
}
