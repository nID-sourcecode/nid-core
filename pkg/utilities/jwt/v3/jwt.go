// Package jwt provides functionality for handling JWT tokens
package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// Duration constants
const (
	DurationDay = time.Hour * 24
)

// error defininitions
var (
	ErrIncorrectBearerLength error = fmt.Errorf("incorrect bearer length")
	ErrInvalidToken          error = fmt.Errorf("invalid token")
)

// ClientOpts options for signin JWT keys
type ClientOpts struct {
	HeaderOpts   HeaderOpts
	SigninMethod jwtgo.SigningMethod
}

// HeaderOpts options for header
type HeaderOpts struct {
	Type string
	Alg  string
	KID  string
}

// DefaultOpts creates new default JWT options
func DefaultOpts() *ClientOpts {
	return &ClientOpts{
		SigninMethod: jwtgo.SigningMethodRS256,
		HeaderOpts: HeaderOpts{
			Type: "JWT",
			Alg:  jwtgo.SigningMethodRS256.Alg(),
			KID:  "1",
		},
	}
}

// NewJWTClient returns a new JWT client
func NewJWTClient(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) *Client {
	return NewJWTClientWithOpts(privKey, pubKey, DefaultOpts())
}

// NewJWTClientWithOpts returns a new JWT client with specified options
func NewJWTClientWithOpts(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, opts *ClientOpts) *Client {
	return &Client{
		PrivKey: privKey,
		PubKey:  pubKey,
		Opts:    opts,
	}
}

// Client client for signin JWT tokens
type Client struct {
	Opts    *ClientOpts
	PrivKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
}

// NewPubKeyClient returns a new JWT public key client
func NewPubKeyClient(pubKey *rsa.PublicKey) *PubKeyClient {
	return NewPubKeyClientWithOpts(pubKey, DefaultOpts())
}

// NewPubKeyClientWithOpts returns a new JWT public key client with specified options
func NewPubKeyClientWithOpts(pubKey *rsa.PublicKey, opts *ClientOpts) *PubKeyClient {
	return &PubKeyClient{
		PubKey: pubKey,
		Opts:   opts,
	}
}

// PubKeyClient client for public key tokens
type PubKeyClient struct {
	Opts   *ClientOpts
	PubKey *rsa.PublicKey
}

// PublicKey returns the public key
func (c *Client) PublicKey() *rsa.PublicKey {
	return c.PubKey
}

// PrivateKey return the priate key
func (c *Client) PrivateKey() *rsa.PrivateKey {
	return c.PrivKey
}

// ValidateAndParseClaims returns the claims for a jwt token.
func (c *Client) ValidateAndParseClaims(bearer string, claims Claims) error {
	_, err := c.validateAndParseToken(bearer, claims)
	if err != nil {
		return err
	}
	return nil
}

// validateAndParseToken parse a jwt string token
func (c *Client) validateAndParseToken(bearer string, claims Claims) (*jwtgo.Token, error) {
	if len(bearer) == 0 {
		return nil, ErrIncorrectBearerLength
	}
	return jwtgo.ParseWithClaims(
		bearer,
		claims,
		func(token *jwtgo.Token) (interface{}, error) {
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwtgo.SigningMethodRSA); !ok {
				return nil, ErrInvalidToken
			}
			// Unpack key from PEM encoded PKCS8
			return c.PublicKey(), nil
		},
	)
}

// SignToken signs a token
func (c *Client) SignToken(claims Claims) (string, error) {
	token := jwtgo.NewWithClaims(c.Opts.SigninMethod, claims)

	token.Header = map[string]interface{}{
		"typ": c.Opts.HeaderOpts.Type,
		"alg": c.Opts.HeaderOpts.Alg,
		"kid": c.Opts.HeaderOpts.KID,
	}

	key := c.PrivateKey()
	if key == nil {
		return "", ErrJWTPrivateNotFound
	}
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
