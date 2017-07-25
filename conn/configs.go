package conn

import (
	"net/url"
)

import (
	"github.com/butlermatt/dslink/crypto"
)

type conf struct {
	isReq  bool
	isResp bool
	broker *url.URL
	name   string
	token  string
	key    *crypto.PrivateKey
}

func IsRequester(c *conf) {
	c.isReq = true
}

func IsResponder(c *conf) {
	c.isResp = true
}

func Broker(brokerUri string) func(c *conf) {
	return func(c *conf) {
		c.broker, _ = url.Parse(brokerUri)
	}
}

func Name(name string) func(c *conf) {
	return func(c *conf) {
		c.name = name
	}
}

func Token(token string) func(c *conf) {
	return func(c *conf) {
		c.token = token
	}
}

func Key(key *crypto.PrivateKey) func(c *conf) {
	return func(c *conf) {
		c.key = key
	}
}
