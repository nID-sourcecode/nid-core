// Package jwt contains functionality for jwt tokens, copies from lab.weave.nl/weave/utilities/jwt
package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

//nolint:gochecknoglobals // the whole package is deprecated anyway
var (
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	keysDirName = "jwtkeys"
	keyFileName = "./" + keysDirName + "/jwt.key"
	pubFileName = "./" + keysDirName + "/jwt.key.pub"
)

func PublicKey() *rsa.PublicKey {
	return publicKey
}

func PrivateKey() *rsa.PrivateKey {
	return privateKey
}

func LoadKeys() error {
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	jwtPub := []byte(os.Getenv("JWT_KEY_PUB"))

	//nolint:nestif
	if len(jwtKey) == 0 || len(jwtPub) == 0 {
		_, err1 := os.Stat(keyFileName)
		_, err2 := os.Stat(pubFileName)
		if os.IsNotExist(err1) || os.IsNotExist(err2) {
			fmt.Println("Warning JWT Keys not found, generating new keypair")
			err := generateKeys()
			if err != nil {
				panic(err)
			}
		} else if err1 != nil {
			return errors.Wrap(err1, "checking for JWT private key file")
		} else if err2 != nil {
			return errors.Wrap(err2, "checking for JWT public key file")
		}
		jwtKey = readFile(keyFileName)
		jwtPub = readFile(pubFileName)
	}

	keyParsed, err := jwtgo.ParseRSAPrivateKeyFromPEM(jwtKey)
	if err != nil {
		return errors.Wrap(err, "loading jwt private key")
	}
	privateKey = keyParsed

	pubParsed, err := jwtgo.ParseRSAPublicKeyFromPEM(jwtPub)
	if err != nil {
		return errors.Wrap(err, "loading public key")
	}
	publicKey = pubParsed

	return nil
}

// LoadPublicKey loads just the jwt public key set in the JWT_KEY_PUB variable, unlike LoadKeys it returns an error if the key cannot be found instead of creating a new one
func LoadPublicKey() error {
	jwtPub := []byte(os.Getenv("JWT_KEY_PUB"))
	if len(jwtPub) == 0 {
		_, err := os.Stat(pubFileName)
		if err != nil {
			return errors.Wrap(err, "finding jwt key")
		}
		jwtPub = readFile(pubFileName)
	}

	pubParsed, err := jwtgo.ParseRSAPublicKeyFromPEM(jwtPub)
	if err != nil {
		return errors.Wrap(err, "loading public key")
	}
	publicKey = pubParsed
	return nil
}

func readFile(fileName string) []byte {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return bytes
}

func generateKeys() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	keyBytes := x509.MarshalPKCS1PrivateKey(key)

	pub := key.PublicKey
	pubBytes, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		panic(err)
	}

	err = savePemFile(keyFileName, "RSA PRIVATE KEY", keyBytes)
	if err != nil {
		panic(err)
	}
	return savePemFile(pubFileName, "RSA PUBLIC KEY", pubBytes)
}

func savePemFile(fileName, pemType string, bytes []byte) error {
	block := &pem.Block{
		Type:  pemType,
		Bytes: bytes,
	}

	newpath := filepath.Join(".", keysDirName)
	os.MkdirAll(newpath, os.ModePerm)

	keyFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer keyFile.Close()

	return pem.Encode(keyFile, block)
}

// ForceFail is method checking mobile session in DB if token contains a mobile session token.
func ForceFail() goa.Middleware {
	forceFail := func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return h(ctx, rw, req)
		}
	}
	fm, _ := goa.NewMiddleware(forceFail)
	return fm
}

func HasScope(scope string, scopes []string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}
