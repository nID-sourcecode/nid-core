// Package keyutil provide utility functionality for keys
package keyutil

import (
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/miscreant/miscreant.go"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

const keyLength int = 32

var errIncorrectDataLength = fmt.Errorf("incorrect data length")

func write(path string, data []byte) error {
	dirname := filepath.Dir(path)
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, os.ModeDir); err != nil {
			return errors.Wrap(err, "unable to create dir for key storage")
		}
	}

	if err := ioutil.WriteFile(path, data, os.ModeAppend); err != nil {
		return errors.Wrap(err, "unable to write key to file")
	}

	return nil
}

func read(path string, length int, generate func() ([]byte, error)) ([]byte, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return generate()
	}

	// nolint: gosec
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read key from path")
	} else if len(data) != length {
		return nil, errors.Wrap(errIncorrectDataLength, "unexpected key length")
	}

	return data, nil
}

func generateKey() ([]byte, error) {
	log.Warn("Key not found, using auto-generated key. This may not be persisted.")
	key := miscreant.GenerateKey(keyLength)
	err := write("keys/key.aes", key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write aes key to file")
	}

	return key, nil
}

// LoadKey load aes keys from file
func LoadKey() ([]byte, error) {
	return read("keys/key.aes", keyLength, generateKey)
}

func generateNonce(c cipher.AEAD) func() ([]byte, error) {
	return func() ([]byte, error) {
		log.Warn("Nonce not found, using auto-generated nonce. This may not be persisted.")
		key := miscreant.GenerateNonce(c)
		err := write(fmt.Sprintf("keys/nonce_%d.bytes", c.NonceSize()), key)

		return key, err
	}
}

// LoadNonce loads nonce from cipher AEAD
func LoadNonce(c cipher.AEAD) ([]byte, error) {
	return read(fmt.Sprintf("keys/nonce_%d.bytes", c.NonceSize()), c.NonceSize(), generateNonce(c))
}
