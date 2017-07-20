package crypto

import (
	"crypto/rand"
	"os"
	"testing"
)

const testKey string = "test.key"

func TestSaveKey(t *testing.T) {
	ed := NewECDH()
	file := defaultKeyFile

	key, err := ed.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	err = SaveKey(key, "")
	if err != nil {
		t.Fatal("Error saving file:", err)
	}

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		t.Fatal("Expected filename does not exist")
	}
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	_ = os.Remove(file)
}

func TestSaveKey2(t *testing.T) {
	ed := NewECDH()
	file := testKey

	key, err := ed.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	err = SaveKey(key, file)
	if err != nil {
		t.Fatalf("Error saving file: %s\nError: %v\n", file, err)
	}

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		t.Fatalf("Expected filename \"%s\" was not created", file)
	}
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	_ = os.Remove(file)
}

func TestLoadKey(t *testing.T) {
	ed := NewECDH()
	file := defaultKeyFile

	key, err := ed.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	err = SaveKey(key, "")
	if err != nil {
		t.Fatalf("Error saving file: %s\nError: %v\n", file, err)
	}

	key2, err := LoadKey("")
	if err != nil {
		t.Fatalf("Error loading key from file: %s\nError: %v\n", file, err)
	}

	cm1, err := ed.Marshal(key)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
	cm2, err := ed.Marshal(key2)
	if err != nil {
		t.Fatal("Unexpected Error", err)
	}
	if cm1 != cm2 {
		t.Fatal("Saved key and loaded key do not match")
	}

	_ = os.Remove(file)
}

func TestLoadKey2(t *testing.T) {
	ed := NewECDH()
	file := testKey

	key, err := ed.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	err = SaveKey(key, file)
	if err != nil {
		t.Fatalf("Error saving file: %s\nError: %v\n", file, err)
	}

	key2, err := LoadKey(file)
	if err != nil {
		t.Fatalf("Error loading key from file: %s\nError: %v\n", file, err)
	}

	cm1, err := ed.Marshal(key)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
	cm2, err := ed.Marshal(key2)
	if err != nil {
		t.Fatal("Unexpected Error", err)
	}
	if cm1 != cm2 {
		t.Fatal("Saved key and loaded key do not match")
	}

	_ = os.Remove(file)
}