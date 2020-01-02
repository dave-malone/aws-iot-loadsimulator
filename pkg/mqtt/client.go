package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client interface {
	pahomqtt.Client
	PublishAsJson(payload interface{}, topic string, qos byte) error
}

type client struct {
	pahomqtt.Client
}

func (c *client) PublishAsJson(payload interface{}, topic string, qos byte) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal json message payload: %v", err)
	}

	if token := c.Publish(topic, qos, false, string(b)); token.Wait() && token.Error() != nil {
		return fmt.Errorf("Failed to publish message: %v\n", token.Error())
	}

	return nil
}

func NewClient(host string, port int, clientID string, tlsConfig *tls.Config) Client {
	clientOpts := &pahomqtt.ClientOptions{
		ClientID:             clientID,
		CleanSession:         true,
		AutoReconnect:        true,
		MaxReconnectInterval: 1 * time.Second,
		KeepAlive:            30,
		TLSConfig:            tlsConfig,
		ProtocolVersion:      3,
	}

	brokerURL := fmt.Sprintf("tcps://%s:%d/mqtt", host, port)
	clientOpts.AddBroker(brokerURL)

	return &client{
		Client: pahomqtt.NewClient(clientOpts),
	}
}

func NewTlsConfig(CertPath, PrivateKeyPath, RootCAPath string) (*tls.Config, error) {
	keypair, err := tls.LoadX509KeyPair(CertPath, PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load X509 keypair: %v", err)
	}

	cert, err := ioutil.ReadFile(RootCAPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load Root CA: %v", err)
	}

	rootCertPool := x509.NewCertPool()
	rootCertPool.AppendCertsFromPEM(cert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{keypair},
		RootCAs:      rootCertPool,
	}

	return tlsConfig, nil
}
