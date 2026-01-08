package pkg

import (
	"net"
	"net/url"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// GetConnectionString return postgres connection string from configs.
func GetConnectionString(cfg DatabaseConfig) string {
	userpass := url.UserPassword(cfg.User, cfg.Password)
	conURL := url.URL{
		Scheme: "postgres",
		User:   userpass,
		Host:   net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:   cfg.Database,
	}

	qs := url.Values{}
	qs.Add("sslmode", cfg.SSLMode)
	conURL.RawQuery = qs.Encode()

	return conURL.String()
}
