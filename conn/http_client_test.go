package conn

import (
	"testing"
	"github.com/butlermatt/dslink/crypto"
	"crypto/rand"
)

func TestNewHttpClient(t *testing.T) {

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
