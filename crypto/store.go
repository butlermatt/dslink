package crypto

import (
	"fmt"
	"io/ioutil"
)

const defaultKeyFile string = ".dslink.key"

// LoadKey will try to load the public and private key configuration from disk.
// If no filename is specified, it will default to .dslink.key.
// This function returns a PrivateKey or error.
func LoadKey(path string) (PrivateKey, error) {
	var priv PrivateKey
	if path == "" {
		path = defaultKeyFile
	}

	d, err := ioutil.ReadFile(path)
	if err != nil {
		return priv, fmt.Errorf("Unable to read file: %s\nError: %v", path, err)
	}

	km := NewECDH()
	return km.Unmarshal(string(d))
}

// SaveKey will attempt to save the specified private key to the specified file.
// If path is not specified, then it will use the default .dslink.key.
// This function will return an error on failure.
func SaveKey(key PrivateKey, path string) error {
	km := NewECDH()

	if path == "" {
		path = defaultKeyFile
	}

	s, err := km.Marshal(key)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, []byte(s), 0644)
}
