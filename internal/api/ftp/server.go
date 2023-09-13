// package ftp создание ftp сервера.
package ftp

import (
	"net/http"
)

// NewServer create ftp server.
func NewServer(address, dir string) *http.Server {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	return &http.Server{
		Addr: address,
	}
}
