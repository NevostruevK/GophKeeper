package client

import (
	"embed"
	"errors"

	//	_ "embed"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

/*
//123go:embed cert/client-cert.pem
var clientCertFile string

//123go:embed cert/client-key.pem
var clientKeyFile string

//123go:embed cert/ca-cert.pem
var clientCACertFile string
*/
//test go:embed cert/ca-cert.pem

//go:embed cert/*
var content embed.FS

//gotst:embed cert/client-key.pem
//var clientKeyFile string

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
	//	fmt.Println(clientCertFile)
	//	fmt.Println(clientKeyFile)
	//	fmt.Println(clientCACertFile)
	fs, err := content.ReadDir("cert")
	fmt.Println(len(fs))
	fmt.Println(fs[0].Name())
	fmt.Println("-------------")
	fileInfo, _ := fs[0].Info()
	fmt.Println("fileinfo- name ------------")
	fmt.Println(fileInfo.Name())
	fmt.Println("----------mode ---")
	fmt.Println(fileInfo.Mode())
	fmt.Println("-------------")
	fmt.Println(fs)
	fmt.Println(err)

	return nil, errors.New("skfsdkfjlksd")
	/*
	   	pemServerCA, err := os.ReadFile(clientCACertFile)
	   //	pemServerCA, err := content.ReadFile("cert/ca-cert.pem")
	   	if err != nil {
	   		return nil, err
	   	}

	   	certPool := x509.NewCertPool()
	   	if !certPool.AppendCertsFromPEM(pemServerCA) {
	   		return nil, fmt.Errorf("failed to add server CA's certificate")
	   	}
	   	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	   	if err != nil {
	   		return nil, err
	   	}

	   	config := &tls.Config{
	   		Certificates: []tls.Certificate{clientCert},
	   		RootCAs:      certPool,
	   	}

	   	return credentials.NewTLS(config), nil
	*/
}
