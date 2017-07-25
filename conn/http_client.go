package conn

import (
	"net/http"
	"net/url"
)

import (
	"fmt"
	"strings"
	"time"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/butlermatt/dslink/log"
	"github.com/butlermatt/dslink/crypto"
)

type dsResp struct {
	Id        string `json:"id"`
	PublicKey string `json:"publicKey"`
	WsUri     string `json:"wsUri"`
	HttpUri   string `json:"httpUri"`
	Version   string `json:"version"`
	TempKey   string `json:"tempKey"`
	Salt      string `json:"salt"`
	SaltS     string `json:"saltS"`
	SaltL     string `json:"saltL"`
	Path      string `json:"path"`
	Format    string `json:"format"`
}

type httpClient struct {
	responder bool
	requester bool
	token     string
	tHash     string
	rawUrl    *url.URL
	privKey   *crypto.PrivateKey
	dsId      string
	keyMaker  crypto.ECDH
	encoder   *Encoder
	htClient  *http.Client
	wsClient  *websocket.Conn
	codecs    map[string]*Encoder
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
		codecs: make(map[string]*Encoder),
	}

	if len(c.token) >= 16 {
		cl.token = c.token[:16]
		cl.tHash = cl.keyMaker.HashToken(cl.dsId, cl.token)
	}

	return cl
}

func (cl *httpClient) Dial() error {
	if len(cl.codecs) <= 0 {
		return fmt.Errorf("no codecs to connect to remote server with")
	}

	resp, err := cl.getWsConfig()
	if err != nil {
		return err
	}

	err = cl.connectWs(resp)

	if err != nil {
		return err
	}

	// TODO: Set timeouts and then handle connections

	return nil
}

func (cl *httpClient) Codec(e *Encoder) {
	cl.codecs[e.Format] = e
}

func (cl *httpClient) getWsConfig() (*dsResp, error) {
	u, _ := url.Parse(cl.rawUrl.String()) // copy url
	q := u.Query()
	q.Add("dsId", cl.dsId)
	if cl.tHash != "" {
		q.Add("token", cl.token+cl.tHash)
	}
	u.RawQuery = q.Encode()

	codecs := make([]string, len(cl.codecs))
	i := 0
	for c := range cl.codecs {
		codecs[i] = "\"" + c + "\""
		i++
	}

	values := fmt.Sprintf("{\"publicKey\": \"%s\", \"isRequester\": %t, \"isResponder\": %t,"+
		"\"linkData\": {}, \"version\": \"1.1.2\", \"formats\": [%s], "+
		"\"enableWebSocketCompression\": true}",
		cl.privKey.PublicKey.Base64(), cl.requester, cl.responder, strings.Join(codecs, ","))

	res, err := cl.htClient.Post(u.String(), "application/json", strings.NewReader(values))
	if err != nil {
		return nil, fmt.Errorf("Error connecting to address: \"%s\"\nError: %s", cl.rawUrl, err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read response: %s", err)
	}

	dr := &dsResp{}
	if err = json.Unmarshal(b, dr); err != nil {
		return nil, fmt.Errorf("Unable to decode response: %s\nError: %s", b, err)
	}
	log.Debug(fmt.Sprintf("Received configuration: %+v\n", *dr))
	return dr, nil
}

func (cl *httpClient) connectWs(conf *dsResp) error {
	cd, ok := cl.codecs[conf.Format]
	if !ok {
		return fmt.Errorf("unknown message encoder %q", conf.Format)
	}

	cl.encoder = cd
	cl.codecs = map[string]*Encoder{} // Clear the others since the server has chosen one.

	pubKey, err := cl.keyMaker.UnmarshalPublic(conf.TempKey)
	if err != nil {
		return fmt.Errorf("unable to parse server key %q, Error: %v", conf.TempKey, err)
	}

	shared := cl.keyMaker.GenerateSharedSecret(*cl.privKey, pubKey)
	auth := cl.keyMaker.HashSalt(conf.Salt, shared)

	u, err := url.Parse(conf.WsUri)
	if err != nil {
		return fmt.Errorf("unable to parse websocket url %q, error: %v", conf.WsUri, err)
	}

	q := u.Query()
	q.Add("auth", auth)
	q.Add("encoder", cl.encoder.Format)
	q.Add("dsId", cl.dsId)
	if cl.tHash != "" {
		q.Add("token", cl.token+cl.tHash)
	}
	u.RawQuery = q.Encode()
	u.Scheme = "ws"

	con, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("unable to connect to websocket at %q, error: %v", conf.WsUri, err)
	}

	cl.wsClient = con
	return nil
}