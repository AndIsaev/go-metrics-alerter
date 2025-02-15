package http

import (
	"crypto/rsa"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/AndIsaev/go-metrics-alerter/internal/service/agent/utils"
)

// Client - define http client
type Client struct {
	Client     *resty.Client
	URL        string
	IPResolver utils.IPResolver
	SecretKey  string
	PublicKey  *rsa.PublicKey
}

// NewClient - func init HTTPClient
func NewClient(url string, ipResolver utils.IPResolver, secretKey string, publicKey *rsa.PublicKey) *Client {
	client := resty.New()
	client.SetTimeout(time.Second * 5)
	return &Client{Client: client, URL: url, IPResolver: ipResolver, SecretKey: secretKey, PublicKey: publicKey}
}
