package tls

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/certifi/gocertifi"
)

func LoadSSLCertificates() (*x509.CertPool, error) {
	certPool, err := gocertifi.CACerts()
	if err != nil {
		return nil, err
	}

	return certPool, nil
}

func CreateTSLClientConfig() (*tls.Config, error) {
	cas, err := LoadSSLCertificates()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs: cas,
	}, nil
}

func CreateClientTransport() (*http.Transport, error) {
	tslClientConfig, err := CreateTSLClientConfig()
	if err != nil {
		return nil, err
	}
	return &http.Transport{
		TLSClientConfig: tslClientConfig,
	}, nil
}
