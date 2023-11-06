// Package jwt provides functionality for handling JWT tokens
package jwt

import (
	"crypto/rsa"
	"fmt"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

// error definitions
var (
	ErrIncorrectBearerLength = fmt.Errorf("incorrect bearer length")
	ErrInvalidToken          = fmt.Errorf("invalid token")
)

// ClientOpts options for signing JWT keys
type ClientOpts struct {
	MarshalSingleStringAsArray bool
	HeaderOpts                 HeaderOpts
	SigningMethod              jwtgo.SigningMethod
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
		MarshalSingleStringAsArray: false,
		SigningMethod:              jwtgo.SigningMethodRS256,
		HeaderOpts: HeaderOpts{
			Type: "JWT",
			Alg:  jwtgo.SigningMethodRS256.Alg(),
			KID:  "1",
		},
	}
}

// NewJWTClientWithOpts returns a new JWT client with specified options
func NewJWTClientWithOpts(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, opts *ClientOpts) *Client {
	return &Client{
		PrivKey: privKey,
		PubKey:  pubKey,
		Opts:    opts,
	}
}

// Client client for signing JWT tokens
type Client struct {
	Opts    *ClientOpts
	PrivKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
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

// PrivateKey return the private key
func (c *Client) PrivateKey() *rsa.PrivateKey {
	return c.PrivKey
}

// ValidateAndParseClaims returns the claims for a jwt token.
func (c *Client) ValidateAndParseClaims(bearer string, claims jwtgo.Claims) error {
	_, err := c.validateAndParseToken(bearer, claims)
	if err != nil {
		return err
	}
	return nil
}

// validateAndParseToken parse a jwt string token
func (c *Client) validateAndParseToken(bearer string, claims jwtgo.Claims) (*jwtgo.Token, error) {
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

// ParseWithClaims validates string token and parses it to a token and claims if validation has passed through.
func (c *Client) ParseWithClaims(token string) (*jwtgo.Token, *DefaultClaims, error) {
	claims := &DefaultClaims{}
	parsedToken, err := jwtgo.ParseWithClaims(token, claims, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodRSA); !ok {
			return nil, ErrInvalidToken
		}

		return c.PrivKey.Public(), nil
	})
	if err != nil {
		return nil, nil, err
	}

	return parsedToken, claims, nil
}

// SignToken signs a token
func (c *Client) SignToken(claims jwtgo.Claims) (string, error) {
	jwtgo.MarshalSingleStringAsArray = c.Opts.MarshalSingleStringAsArray
	token := jwtgo.NewWithClaims(c.Opts.SigningMethod, claims)

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
