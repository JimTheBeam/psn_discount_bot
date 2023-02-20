package httpclient

import (
	"errors"
	"time"
)

type Config struct {
	Timeout             time.Duration `yaml:"timeout"`
	TLSHandshakeTimeout time.Duration `yaml:"tls_handshake_timeout"`
	DialerTimeout       time.Duration `yaml:"dialer_timeout"`
}

func (cfg *Config) Validate() error {
	if cfg.Timeout == 0 {
		return errors.New("timeout is empty")
	}

	if cfg.TLSHandshakeTimeout == 0 {
		return errors.New("tls_handshake_timeout is empty")
	}

	if cfg.DialerTimeout == 0 {
		return errors.New("dialer_timeout is empty")
	}

	return nil
}
