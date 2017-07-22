package conn

import (
	"net/url"
	"net/http"
)

import (
	"github.com/butlermatt/dslink/crypto"
	"time"
)

type httpClient struct {
	responder bool
	requester bool
	token     string
	tHash     string
	rawUrl    *url.URL
	privKey   *crypto.PrivateKey
	dsId      string
	keyMaker  crypto.ECDH
	htClient  *http.Client
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
		responder: c.isResp,
		requester: c.isReq,
		rawUrl:    c.broker,
		privKey:   c.key,
		dsId:      c.key.DsId(c.name),
		keyMaker:  crypto.NewECDH(),
		htClient:  &http.Client{Timeout: time.Minute},
	}

	if len(c.token) >= 16 {
		cl.token = c.token[:16]
		cl.tHash = cl.keyMaker.HashToken(cl.dsId, cl.token)
	}

	return cl
}
