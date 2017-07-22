package conn

import (
	"testing"
	"github.com/butlermatt/dslink/crypto"
	"crypto/rand"
	"net/url"
)

func TestNewHttpClient(t *testing.T) {
	a := "http://localhost:8080/conn"
	n := "test"
	tok := "12345678901234567891"
	km := crypto.NewECDH()
	key, err := km.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("Error generating private key")
	}

	cl := NewHttpClient(IsResponder, Name(n), Token(tok), Broker(a), Key(&key))

	addr, _ := url.Parse(a)
	if addr.String() != cl.rawUrl.String() {
		t.Errorf("httpClient.rawUrl did not match. expected=%q got=%q", addr, cl.rawUrl)
	}

	if !cl.responder {
		t.Errorf("httpClient.responder was not set. expected=%q got=%q", true, cl.responder)
	}

	if cl.requester {
		t.Errorf("httpClient.requester should be false. got=%q", cl.requester)
	}

	if cl.privKey != &key {
		t.Error("httpClient.privKey does not match expected private key")
	}

	dsid := key.DsId(n)
	if cl.dsId != dsid {
		t.Errorf("httpClient.dsId did not match. expected=%q got=%q", dsid, cl.dsId)
	}

	if cl.token != tok[:16] {
		t.Errorf("httpClient.token does not match. expected=%q got=%q", tok[:16], cl.token)
	}

	thash := km.HashToken(dsid, cl.token)
	if thash != cl.tHash {
		t.Errorf("httpClient.tHash does not match. expected=%q got=%q", thash, cl.tHash)
	}

	if cl.htClient == nil {
		t.Error("httpClient.htClient should not be nil")
	}
}

func TestNewHttpClientPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Did not panic when expected")
		}
	}()

	// No Key
	_ = NewHttpClient(Broker("http://localhost:8080/conn"), Name("Test-"))

}

func TestNewHttpClientPanic2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Did not panic when expected")
		}
	}()

	km := crypto.NewECDH()
	pk, _ := km.GenerateKey(rand.Reader)
	// No Name
	_ = NewHttpClient(Broker("http://localhost:8080/conn"), Key(&pk))
}

func TestNewHttpClientPanic3(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Did not panic when expected")
		}
	}()

	km := crypto.NewECDH()
	pk, _ := km.GenerateKey(rand.Reader)
	// No Broker
	_ = NewHttpClient(Name("Test-"), Key(&pk))
}
