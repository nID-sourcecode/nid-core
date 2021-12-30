// Package jwtconfig provides the jwt config
package jwtconfig

import (
	"crypto/rsa"
	"io/ioutil"
	"path/filepath"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v2"
)

// JWTKey contains the secret parameters for jwt
type JWTKey struct {
	*rsa.PrivateKey
	ID string
}

// Read reads the jwt secrets from a directory
func Read(path string) (*JWTKey, error) {
	kid, err := ioutil.ReadFile(filepath.Clean(path + "/kid"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to read kid from file")
	}
	keyBytes, err := ioutil.ReadFile(filepath.Clean(path + "/jwtkey.pem"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to read JWT key from file")
	}

	privKey, err := jwt.ParseKey(keyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse JWT key bytes")
	}

	return &JWTKey{
		PrivateKey: privKey,
		ID:         string(kid),
	}, nil
}
