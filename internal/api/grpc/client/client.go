package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const refreshDuration = time.Hour

type Client struct {
	Auth   *AuthClient
	Keeper *KeeperClient
}

func NewClient(address string, enableTLS bool) (*Client, error) {
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
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
		log.Println("cannot dial server: ", err)
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

func (c *Client) Close() error {
	errAuth := c.Auth.Close()
	errKeeper := c.Keeper.Close()
	if errAuth != nil || errKeeper != nil {
		return fmt.Errorf("gRPC client closed with err: auth: %v , keeper: %v", errAuth, errKeeper)
	}
	return nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	const (
		clientCertFile   = "../../cert/client-cert.pem"
		clientKeyFile    = "../../cert/client-key.pem"
		clientCACertFile = "../../cert/ca-cert.pem"
	)
	pemServerCA, err := os.ReadFile(clientCACertFile)
	if err != nil {
		log.Println("0:", err)
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		log.Println("1:", err)
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Println("2:", err)
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}
