package config

import "time"

// Config contains program's options.
type Config struct {
	TokenDuration time.Duration `json:"token_duration"`
	Address       string        `json:"address"`
	FtpAddress    string        `json:"ftp_address"`
	FtpDir        string        `json:"ftp_dir"`
	TokenKey      string        `json:"token_key"`
	DSN           string        `json:"dsn"`
	EnableTLS     bool          `json:"use_tls"`
	CryptoKey     string        `json:"crypto_key"`
	CryptoNonce   string        `json:"crypto_nonce"`
}
