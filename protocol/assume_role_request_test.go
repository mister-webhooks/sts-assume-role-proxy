package protocol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleCredentials(t *testing.T) {
	t.Run("encoding round-trips", func(t *testing.T) {
		rc := RoleCredentials{
			Result:          0x10,
			Expiration:      time.Now(),
			AccessKeyId:     "access",
			SecretAccessKey: "secret",
			SessionToken:    "token",
		}
		encoded, err := rc.MarshalBinary()
		decoded := new(RoleCredentials)

		require.NoError(t, err)

		err = decoded.UnmarshalBinary(encoded)

		assert.NoError(t, err)
		assert.Equal(t, rc.Result, decoded.Result)
		assert.True(t, decoded.Expiration.Equal(rc.Expiration))
		assert.Equal(t, rc.AccessKeyId, decoded.AccessKeyId)
		assert.Equal(t, rc.SecretAccessKey, decoded.SecretAccessKey)
		assert.Equal(t, rc.SessionToken, decoded.SessionToken)
	})
}
