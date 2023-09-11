// package client gRPC client
package client

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed cert/*
var content embed.FS

const refreshDuration = time.Hour

// Client gRPC client.
type Client struct {
	Auth   *AuthClient   // клиент для авторизации пользователя.
	Keeper *KeeperClient // клиент для обмена данными.
}

// NewClient returns a new gRPC client
func NewClient(address string, enableTLS bool) (*Client, error) {
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			return nil, fmt.Errorf("cannot load TLS credentials: %v", err)
		}

		transportOption = grpc.WithTransportCredentials(tlsCredentials)
	}

	conn, err := grpc.Dial(address, transportOption)
	if err != nil {
		return nil, fmt.Errorf("cannot dial server: %v", err)
	}
	interceptor := NewAuthInterceptor(authMethods(), refreshDuration)

	authClient := NewAuthClient(conn, interceptor)

	conn, err = grpc.Dial(
		address,
		transportOption,
		grpc.WithUnaryInterceptor(interceptor.Unary()),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot dial server: %v", err)
	}

	KeeperClient := NewKeeperClient(conn)

	return &Client{authClient, KeeperClient}, nil
}

func authMethods() map[string]bool {
	const keeperServicePath = "/gophKeeper.Keeper/"

	return map[string]bool{
		keeperServicePath + "GetSpecs":       true,
		keeperServicePath + "GetSpecsOfType": true,
		keeperServicePath + "GetData":        true,
		keeperServicePath + "GetDescription": true,
		keeperServicePath + "AddRecord":      true,
	}
}

// Close освобождение ресурсов gRPC клиента.
func (c *Client) Close() error {
	errAuth := c.Auth.Close()
	errKeeper := c.Keeper.Close()
	if errAuth != nil || errKeeper != nil {
		return fmt.Errorf("gRPC client closed with err: auth: %v , keeper: %v", errAuth, errKeeper)
	}
	return nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	/*	const (
			clientCertFile   = "../../cert/client-cert.pem"
			clientKeyFile    = "../../cert/client-key.pem"
			clientCACertFile = "../../cert/ca-cert.pem"
		)
	*/
	const (
		clientCertFile   = "cert/client-cert.pem"
		clientKeyFile    = "cert/client-key.pem"
		clientCACertFile = "cert/ca-cert.pem"
	)

	//	pemServerCA, err := os.ReadFile(clientCACertFile)
	pemServerCA, err := content.ReadFile(clientCACertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	certPEMBlock, err := content.ReadFile(clientCertFile)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := content.ReadFile(clientKeyFile)
	if err != nil {
		return nil, err
	}

	//	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	clientCert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)

	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil

}
