package tls

import (
	"crypto/tls"

	"google.golang.org/grpc/credentials"
)

func GetTransportCredentials(tlsCert string, tlsKey string) credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		panic(err)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	return creds
}
