package conn

import (
	"testing"
	"net/url"
	"crypto/rand"
	"github.com/butlermatt/dslink/crypto"
)

func TestIsRequester(t *testing.T) {
	c := new(conf)

	applyTest(c, IsRequester)
	if !c.isReq {
		t.Errorf("conf.isReq expected=%v got=%v", true, c.isReq)
	}
}

func TestIsResponder(t *testing.T) {
	c := new(conf)

	applyTest(c, IsResponder)
	if !c.isResp {
		t.Errorf("conf.isResp expect=%v got=%v", true, c.isResp)
	}
}

func TestName(t *testing.T) {
	c := new(conf)
	n := "test-"

	applyTest(c, Name(n))
	if c.name != n {
		t.Errorf("conf.name expected=%q got=%q", n, c.name)
	}
}

func TestToken(t *testing.T) {
	c := new(conf)
	tok := "abcdefghijklmnopqrstuvwxyz"

	applyTest(c, Token(tok))
	if c.token != tok {
		t.Errorf("conf.token expect=%q got=%q", tok, c.token)
	}
}

func TestBroker(t *testing.T) {
	c := new(conf)
	a := "http://localhost8080/conn"

	applyTest(c, Broker(a))
	tmp, _ := url.Parse(a)
	if tmp.String() != c.broker.String() {
		t.Errorf("conf.Broker expected=%v got=%v", tmp, c.broker)
	}
}

func TestKey(t *testing.T) {
	c := new(conf)
	km := crypto.NewECDH()
	key, err := km.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("Failed to generate key for test", err)
	}

	applyTest(c, Key(&key))
	if c.key != &key {
		t.Errorf("conf.key does not match generated key")
	}
}

func TestMultiple(t *testing.T) {
	c := new(conf)
	tok := "abc"
	n := "test-"
	u := "http://localhost:8080/conn"
	uri, _ := url.Parse(u)
	km := crypto.NewECDH()
	key, err := km.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("Failed to generate private key for test", err)
	}

	applyTest(c, IsResponder, IsRequester, Name(n), Token(tok), Broker(u), Key(&key))
	if !c.isResp && !c.isReq && c.name != n && c.token != tok &&
		 c.broker.String() != uri.String() && c.key != &key {
		t.Errorf("Failed to update all values.\n" +
		  "c.isResp expected=%v got=%v\n" +
		  "c.isReq  expected=%v got=%v\n" +
		  "c.name   expected=%q got=%q\n" +
		  "c.token  expected=%q got=%q\n" +
		  "c.broker expected=%q got=%q\n" +
		  "c.key    expected=%+v got=%+v",
		true, c.isResp,
		true, c.isReq,
		n, c.name,
		tok, c.token,
		uri, c.broker,
		key, c.key)
	}
}

func applyTest(c *conf, opts ...func(c *conf)) {
	for _, o := range opts {
		o(c)
	}
}