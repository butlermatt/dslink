package conn

import (
	"net/url"
)

import (
	"github.com/butlermatt/dslink/crypto"
)

type httpClient struct {
	keyMaker  crypto.ECDH
	privKey   *crypto.PrivateKey
	rawUrl    *url.URL
	dsId      string
	responder bool
	requester bool
}

func NewHttpClient(opts ...func(c *conf)) *httpClient {
	c := &conf{}

	for _, opt := range opts {
		opt(c)
	}

	if c.broker == nil {
		panic("cannot create httpClient without broker address")
	}

	if c.name == "" {
		panic("cannot create httpClient without link name")
	}

	if c.key == nil {
		panic("cannot create httpClient without a private key")
	}

	cl := &httpClient{
		keyMaker:  crypto.NewECDH(),
		privKey:   c.key,
		rawUrl:    c.broker,
		responder: c.isResp,
		requester: c.isReq,
		dsId:      c.key.DsId(c.name),
	}

	return cl
}
