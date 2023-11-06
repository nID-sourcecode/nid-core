// Package keymanager provides key manager for managing keys
package keymanager

import (
	"context"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

const byteLength = 8

// Error definitions
var (
	errNoRSAKeyFound = fmt.Errorf("no rsa key for encryption found in jwks")
)

// IJWKSFetcher interface for fetching JWKS
type IJWKSFetcher interface {
	Fetch(string) (jwk.Set, error)
}

// JWKSFetcher implementation of the IJWKSFetcher interface
type JWKSFetcher struct{}

// Fetch fetches as jwk url
func (j *JWKSFetcher) Fetch(url string) (jwk.Set, error) {
	return jwk.Fetch(context.TODO(), url)
}

// KeyManager manages keys
type KeyManager struct {
	cache          *ttlcache.Cache
	jwkURLTemplate string // e.g. "pseudo.{{namespace}}/jwks"
	jwksFetcher    IJWKSFetcher
}

// NewKeyManager creates new key manager
func NewKeyManager(jwkURLTemplate string, cacheDuration time.Duration, fetcher IJWKSFetcher) KeyManager {
	cache := ttlcache.NewCache()
	cache.SetTTL(cacheDuration)

	return KeyManager{
		jwksFetcher:    fetcher,
		cache:          cache,
		jwkURLTemplate: jwkURLTemplate,
	}
}

// GetKey retrieves key for given namespace
func (k KeyManager) GetKey(namespace string) (*rsa.PublicKey, error) {
	if key, ok := k.cache.Get(namespace); ok {
		if rsaKey, ok := key.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}
		k.cache.Remove(namespace)
	}

	jwksURL := strings.ReplaceAll(k.jwkURLTemplate, "{{namespace}}", namespace)
	set, err := k.jwksFetcher.Fetch(jwksURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch jwks")
	}

	keyIterator := set.Keys(context.TODO())
	for keyIterator.Next(context.TODO()) {
		value := keyIterator.Pair().Value
		key := value.(jwk.Key)
		if key.Algorithm() == jwa.RSA1_5 && key.KeyUsage() == string(jwk.ForEncryption) {
			rsaKey := parseRsaPublicKeyFromJwk(key)

			k.cache.Set(namespace, rsaKey)

			return rsaKey, nil
		}
	}

	return nil, errors.Wrap(errNoRSAKeyFound, "unable to find rsa key for given namespace")
}

// Cleanup cleansup key manager cache
func (k KeyManager) Cleanup() {
	k.cache.Close()
}

func parseRsaPublicKeyFromJwk(key jwk.Key) *rsa.PublicKey {
	rsaPublicKey := key.(jwk.RSAPublicKey)
	nBytes := rsaPublicKey.N()

	eBytes := rsaPublicKey.E()
	ePaddedBytes := make([]byte, byteLength-len(eBytes), byteLength)
	ePaddedBytes = append(ePaddedBytes, eBytes...)
	e := binary.BigEndian.Uint64(ePaddedBytes)

	n := big.Int{}
	n.SetBytes(nBytes)

	return &rsa.PublicKey{
		N: &n,
		E: int(e),
	}
}
