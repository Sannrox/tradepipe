package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func WaitForPortToBeNotAttachedWithLimit(port string, limit int) error {
	for i := 0; i < limit; i++ {
		l, err := net.Listen("tcp", "localhost:"+port)
		if err == nil {
			l.Close()
			return nil
		}
		logrus.Print(err)
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout waiting for port %s to be not attached to the service", port)
}

func WaitForRestServerToBeUp(url string, limit int) error {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

	for i := 0; i < limit; i++ {
		_, err := client.Get(url)
		if err == nil {
			return nil
		}
		logrus.Print(err)
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout waiting for rest server to be up")
}
