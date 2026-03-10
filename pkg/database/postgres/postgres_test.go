package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnectionString(t *testing.T) {
	t.Run("ssl_mode_disabled", func(t *testing.T) {
		cfg := &DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "dev",
			Password: "psswd",
			Database: "control-panel",
			SSLMode:  "disable",
		}

		str := GetConnectionString(cfg)
		assert.Equal(t, str, "postgres://dev:psswd@localhost:5432/control-panel?sslmode=disable")
	})
	t.Run("ssl_mode_allow", func(t *testing.T) {
		cfg := &DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "dev",
			Password: "psswd",
			Database: "control-panel",
			SSLMode:  "allow",
		}

		str := GetConnectionString(cfg)
		assert.Equal(t, str, "postgres://dev:psswd@localhost:5432/control-panel?sslmode=allow")
	})
	t.Run("complicate_password", func(t *testing.T) {
		cfg := &DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "dev",
			Password: "!@#$%^&*()UYTGHJNDyvghsdbu",
			Database: "control-panel",
			SSLMode:  "allow",
		}

		str := GetConnectionString(cfg)
		assert.Equal(t, str, "postgres://dev:%21%40%23$%25%5E&%2A%28%29UYTGHJNDyvghsdbu@localhost:5432/control-panel?sslmode=allow")
	})
}
