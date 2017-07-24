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

type msgFormat int

const (
	fmtJson msgFormat = iota
	fmtMsgP
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
	format    msgFormat
	htClient  *http.Client
	wsClient  *websocket.Conn
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

func (cl *httpClient) Dial() error {
	resp, err := cl.getWsConfig()
	if err != nil {
		return err
	}

	err = cl.connectWs(resp)

	if err != nil {
		return err
	}

	// Set timeouts and then handle connections

	return nil
}

func (cl *httpClient) getWsConfig() (*dsResp, error) {
	u, _ := url.Parse(cl.rawUrl.String()) // copy url
	q := u.Query()
	q.Add("dsId", cl.dsId)
	if cl.tHash != "" {
		q.Add("token", cl.token+cl.tHash)
	}
	u.RawQuery = q.Encode()

	values := fmt.Sprintf("{\"publicKey\": \"%s\", \"isRequester\": %t, \"isResponder\": %t,"+
		"\"linkData\": {}, \"version\": \"1.1.2\", \"formats\": [\"msgpack\",\"json\"], "+
		"\"enableWebSocketCompression\": true}",
		cl.privKey.PublicKey.Base64(), cl.requester, cl.responder)

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
	switch conf.Format {
	case "json":
		cl.format = fmtJson
	case "msgpack":
		cl.format = fmtMsgP
	default:
		return fmt.Errorf("unknown message format %q", conf.Format)
	}

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
	q.Add("format", conf.Format)
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