package config_test

import (
	"testing"

	"github.com/mister-webhooks/sts-assume-role-proxy/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigurationFromYAML(t *testing.T) {
	t.Run("it deserializes a correctly-formed server configuration", func(t *testing.T) {
		sc, err := config.NewConfigurationFromYAML([]byte(`
rules:
  '[root]':
    123: 'arn:example:foo'
    456: 'arn:example:bar'
  ns1:
    456: 'arn:example:qux'
    789: 'arn:example:quux'
    `))

		require.NoError(t, err)

		root123, ok := sc.AccessTable.Lookup("[root]", 123)
		require.True(t, ok)
		assert.Equal(t, "arn:example:foo", root123)
	})

	t.Run("it deserializes a server configuration with an empty ruleset", func(t *testing.T) {
		_, err := config.NewConfigurationFromYAML([]byte(`rules:
    `))

		require.NoError(t, err)
	})
}
