// package config options for client and server.
package config

import "time"

// Config contains program's options.
type Config struct {
	TokenDuration time.Duration `json:"token_duration"` // token duration
	Address       string        `json:"address"`        // gRPC address
	FtpAddress    string        `json:"ftp_address"`    // FPT server address
	FtpDir        string        `json:"ftp_dir"`        // FTP directory
	TokenKey      string        `json:"token_key"`      // token key
	DSN           string        `json:"dsn"`            // dsn
	EnableTLS     bool          `json:"use_tls"`        // enable gRPC TLS
	CryptoKey     string        `json:"crypto_key"`     // crypto key
	CryptoNonce   string        `json:"crypto_nonce"`   // crypto nonce
}
