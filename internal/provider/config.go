package provider

import (
	"errors"
)

// Config contains all available configuration options.
type Config struct {
	APIEndpoint string
	APIToken    string
}

// Validate performs config validation.
func (c *Config) Client() error {
	if c.APIToken == "" {
		return errors.New("API_TOKEN must be not empty")
	}
	if c.APIEndpoint == "" {
		c.APIEndpoint = "api.servicepipe.ru/api/v1"
	}

	// tflog.Info(ctx, fmt.Sprintf("cloudflare Client configured: %s", c.Email))

	return nil
}
