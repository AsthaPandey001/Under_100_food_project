package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	_, err := GetConfig()
	require.NoError(t, err)

	cfg, err := GetConfig()
	require.NoError(t, err)
	assert.Equal(t, EnvLocal, cfg.App.Environment)
}

func TestConfig_Validate(t *testing.T) {
	t.Run("Valid local environment", func(t *testing.T) {
		cfg := &Config{
			App: AppConfig{
				Environment: EnvLocal,
			},
		}
		assert.NoError(t, cfg.Validate())
	})

	t.Run("Production fails with default encryption key", func(t *testing.T) {
		cfg := &Config{
			App: AppConfig{
				Environment:   EnvProduction,
				EncryptionKey: "?E(G+KbPeShVmYq3t6w9z$C&F)J@McQf",
			},
			Database: DatabaseConfig{
				Password: "strong-db-password",
			},
			HTTP: HTTPConfig{
				Domain: "bite100.example.com",
			},
		}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PAGODA_APP_ENCRYPTIONKEY must be set")
	})

	t.Run("Production passes with valid production settings", func(t *testing.T) {
		cfg := &Config{
			App: AppConfig{
				Environment:   EnvProduction,
				EncryptionKey: "a-very-secure-32-byte-production-key!",
			},
			Database: DatabaseConfig{
				Password: "super-strong-db-password",
			},
			HTTP: HTTPConfig{
				Domain: "bite100.example.com",
			},
		}
		assert.NoError(t, cfg.Validate())
	})
}

