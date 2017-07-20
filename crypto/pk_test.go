package crypto

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	clientPrivate     = "M6S41GAL0gH0I97Hhy7A2-icf8dHnxXPmYIRwem03HE"
	clientPublic      = "BEACGownMzthVjNFT7Ry-RPX395kPSoUqhQ_H_vz0dZzs5RYoVJKA16XZhdYd__ksJP0DOlwQXAvoDjSMWAhkg4"
	clientDsId        = "test-s-R9RKdvC2VNkfRwpNDMMpmT_YWVbhPLfbIc-7g4cpc"
	serverTempPrivate = "rL23cF6HxmEoIaR0V2aORlQVq2LLn20FCi4_lNdeRkk"
	serverTempPublic  = "BCVrEhPXmozrKAextseekQauwrRz3lz2sj56td9j09Oajar0RoVR5Uo95AVuuws1vVEbDzhOUu7freU0BXD759U"
	sharedSecret      = "116128c016cf380933c4b40ffeee8ef5999167f5c3d49298ba2ebfd0502e74e3"
	hashedAuth        = "V2P1nwhoENIi7SqkNBuRFcoc8daWd_iWYYDh_0Z01rs"
)

func TestNewECDH(t *testing.T) {
	ecdh := NewECDH()

	if ecdh == nil {
		t.Fatal("ecdh was nil!")
	}
}

func TestEllipticECDH_Unmarshal(t *testing.T) {
	ecdh := NewECDH()

	_, err := ecdh.Unmarshal(clientPrivate)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	_, err = ecdh.Unmarshal(serverTempPrivate)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
}

func TestEllipticECDH_Unmarshal2(t *testing.T) {
	ecdh := NewECDH()

	tmpC := fmt.Sprintf("%s %s", clientPrivate, clientPublic)
	priv1, err := ecdh.Unmarshal(tmpC)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	priv2, err := ecdh.Unmarshal(clientPrivate)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if priv1.Base64() != priv2.Base64() {
		t.Errorf("D mismatch!\n%X\n%X", priv1.Base64(), priv2.Base64())
	}

}

func TestEllipticECDH_UnmarshalPublic(t *testing.T) {
	ecdh := NewECDH()

	_, err := ecdh.UnmarshalPublic(clientPublic)
	if err != nil {
		t.Errorf("UnmarshalPublic failed: %v", err)
	}

	_, err = ecdh.UnmarshalPublic(serverTempPublic)
	if err != nil {
		t.Errorf("UnmarshalPublic failed: %v", err)
	}
}

func TestPublicKey_Base64(t *testing.T) {
	ecdh := NewECDH()
	cpriv, _ := ecdh.Unmarshal(clientPrivate)
	if cpriv.PublicKey.Base64() != clientPublic {
		t.Errorf("%v != %v", cpriv.PublicKey.Base64(), clientPublic)
	}

	cpub, _ := ecdh.UnmarshalPublic(clientPublic)
	if cpub.Base64() != clientPublic {
		t.Errorf("%v != %v", cpub.Base64(), clientPublic)
	}

	if cpriv.PublicKey.Base64() != cpub.Base64() {
		t.Errorf("%v != %v", cpriv.Base64(), cpub.Base64())
	}

	spriv, _ := ecdh.Unmarshal(serverTempPrivate)
	if spriv.PublicKey.Base64() != serverTempPublic {
		t.Errorf("%v != %v", cpriv.PublicKey.Base64(), clientPublic)
	}

	spub, _ := ecdh.UnmarshalPublic(serverTempPublic)
	if spub.Base64() != serverTempPublic {
		t.Errorf("%v != %v", spub.Base64(), clientPublic)
	}

	if spriv.PublicKey.Base64() != spub.Base64() {
		t.Errorf("%v != %v", spriv.Base64(), spub.Base64())
	}
}

func TestPublicKey_DsId(t *testing.T) {
	ecdh := NewECDH()
	priv, _ := ecdh.Unmarshal(clientPrivate)
	if priv.PublicKey.DsId("test-") != clientDsId {
		t.Errorf("%v != %v", priv.PublicKey.DsId("test-"), clientDsId)
	}
}

func TestPublicKey_VerifyDsId(t *testing.T) {
	ecdh := NewECDH()
	priv, _ := ecdh.Unmarshal(clientPrivate)
	if !priv.PublicKey.VerifyDsId(clientDsId) {
		t.Errorf("VerifyDsId failed: %s did not verify with %s", clientDsId, priv.Base64())
	}
}

func TestPublicKey_Hash64(t *testing.T) {
	ecdh := NewECDH()
	priv, _ := ecdh.Unmarshal(clientPrivate)
	pub, _ := ecdh.UnmarshalPublic(clientPublic)

	if priv.Hash64() != pub.Hash64() {
		t.Errorf("%v != %v", priv.Hash64(), pub.Hash64())
	}

	priv, _ = ecdh.Unmarshal(serverTempPrivate)
	pub, _ = ecdh.UnmarshalPublic(serverTempPublic)

	if priv.Hash64() != pub.Hash64() {
		t.Errorf("%v != %v", priv.Hash64(), pub.Hash64())
	}
}

func TestEllipticECDH_Marshal(t *testing.T) {
	ecdh := NewECDH()
	priv, _ := ecdh.Unmarshal(clientPrivate)

	tmp := fmt.Sprintf("%s %s", clientPrivate, clientPublic)
	m, _ := ecdh.Marshal(priv)
	if tmp != m {
		t.Errorf("marshal mismatch:\n%s\n%s", tmp, m)
	}

	priv, _ = ecdh.Unmarshal(serverTempPrivate)
	tmp = fmt.Sprintf("%s %s", serverTempPrivate, serverTempPublic)
	m, _ = ecdh.Marshal(priv)
	if tmp != m {
		t.Errorf("marshal mismatch:\n%s\n%s", tmp, m)
	}
}

func TestEllipticECDH_GenerateSharedSecret(t *testing.T) {
	ecdh := NewECDH()

	cpriv, _ := ecdh.Unmarshal(clientPrivate)
	spriv, _ := ecdh.Unmarshal(serverTempPrivate)

	cshare := ecdh.GenerateSharedSecret(cpriv, spriv.PublicKey)
	sshare := ecdh.GenerateSharedSecret(spriv, cpriv.PublicKey)

	if hex.EncodeToString(cshare) != sharedSecret {
		t.Errorf("%v != %v", hex.EncodeToString(cshare), sharedSecret)
	}

	if hex.EncodeToString(sshare) != sharedSecret {
		t.Errorf("%v != %v", hex.EncodeToString(sshare), sharedSecret)
	}

	if !bytes.Equal(cshare, sshare) {
		t.Errorf("%v != %v", hex.EncodeToString(cshare), hex.EncodeToString(sshare))
	}
}

func TestEllipticECDH_HashSalt(t *testing.T) {
	ecdh := NewECDH()

	spriv, _ := ecdh.Unmarshal(serverTempPrivate)
	cpriv, _ := ecdh.Unmarshal(clientPrivate)

	share := ecdh.GenerateSharedSecret(spriv, cpriv.PublicKey)
	hs := ecdh.HashSalt("0000", share)
	if hs != hashedAuth {
		t.Errorf("%v != %v", hs, hashedAuth)
	}

	share = ecdh.GenerateSharedSecret(cpriv, spriv.PublicKey)
	hs = ecdh.HashSalt("0000", share)
	if hs != hashedAuth {
		t.Errorf("%v != %v", hs, hashedAuth)
	}
}

func TestEllipticECDH_GenerateKey(t *testing.T) {
	ecdh := NewECDH()

	cpriv, err := ecdh.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	cpriv2, err := ecdh.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	cm, err := ecdh.Marshal(cpriv)
	if err != nil {
		t.Fatalf("failed to marshal private key, %v", err)
	}
	cm2, err := ecdh.Marshal(cpriv2)
	if err != nil {
		t.Fatalf("failed to marshal private key, %v", err)
	}
	tmp := fmt.Sprintf("%s %s", clientPrivate, clientPublic)
	if cm == tmp {
		t.Error("Generated key matches static known values.")
	}
	tmp = fmt.Sprintf("%s %s", serverTempPrivate, serverTempPublic)
	if cm == tmp {
		t.Error("Generated key matches static known values.")
	}
	if cm == cm2 {
		t.Errorf("2 Generated Keys match: %v", cm)
	}

	spriv, _ := ecdh.Unmarshal(serverTempPrivate)

	cshare := ecdh.GenerateSharedSecret(cpriv, spriv.PublicKey)
	sshare := ecdh.GenerateSharedSecret(spriv, cpriv.PublicKey)
	if !bytes.Equal(cshare, sshare) {
		t.Errorf("Failed to generate equal shared secrets:\n%v\n%v", hex.EncodeToString(cshare), hex.EncodeToString(sshare))
	}
}
