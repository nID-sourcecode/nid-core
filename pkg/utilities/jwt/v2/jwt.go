// Package jwt provides functionality for handling JWT tokens
package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
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
	ClaimsOpts   ClaimsOpts
	SigninMethod jwtgo.SigningMethod
}

// HeaderOpts options for header
type HeaderOpts struct {
	Type string
	Alg  string
	KID  string
}

// ClaimsOpts options for claims
type ClaimsOpts struct {
	Issuer     string
	Audience   []string
	Expiration ExpirationFunc
	JWTID      IDFunc
	IssuedAt   IssuedAtFunc
	NotBefore  NotBeforeFunc
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
		ClaimsOpts: ClaimsOpts{
			Issuer:     "weave.nl",
			Audience:   []string{"weave"},
			Expiration: expiration,
			JWTID:      jwtID,
			IssuedAt:   issuedAt,
			NotBefore:  notBefore,
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

// ExpirationFunc function definition for expiration
type ExpirationFunc func() int64

func expiration() int64 {
	return time.Now().Add(DurationDay).Unix()
}

// IDFunc function definition for creating JWT ID
type IDFunc func() string

func jwtID() string {
	return uuid.Must(uuid.NewV4()).String()
}

// IssuedAtFunc function definition for when the JWT was created
type IssuedAtFunc func() int64

func issuedAt() int64 {
	return time.Now().Unix()
}

// NotBeforeFunc function definition for creating the time before the JWT may not be used
type NotBeforeFunc func() int64

func notBefore() int64 {
	return time.Now().Add(-2 * time.Minute).Unix()
}

// PublicKey returns the public key
func (c *Client) PublicKey() *rsa.PublicKey {
	return c.PubKey
}

// PrivateKey return the priate key
func (c *Client) PrivateKey() *rsa.PrivateKey {
	return c.PrivKey
}

// GetScopesFromClaims retrieve scopes from claims
func (c *Client) GetScopesFromClaims(claims jwtgo.MapClaims) ([]string, error) {
	scopes, ok := claims["scope"].([]interface{})
	if !ok {
		return nil, ErrInvalidToken
	}
	res := []string{}
	for _, scope := range scopes {
		scopeString, ok := scope.(string)
		if !ok {
			return nil, ErrInvalidToken
		}
		res = append(res, scopeString)
	}
	return res, nil
}

// GetClaims returns the claims for a jwt token.
func (c *Client) GetClaims(bearer string) (jwtgo.MapClaims, error) {
	token, err := c.parseToken(bearer)
	if err != nil {
		return nil, err
	}
	return c.getClaims(token), nil
}

// getClaims retrieve claims from token
func (c *Client) getClaims(token *jwtgo.Token) jwtgo.MapClaims {
	return token.Claims.(jwtgo.MapClaims)
}

// parseToken parse a jwt string token
func (c *Client) parseToken(bearer string) (*jwtgo.Token, error) {
	if len(bearer) == 0 {
		return nil, ErrIncorrectBearerLength
	}
	return jwtgo.ParseWithClaims(
		bearer,
		jwtgo.MapClaims{},
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
func (c *Client) SignToken(customClaims map[string]interface{}) (string, error) {
	token := jwtgo.New(c.Opts.SigninMethod)

	token.Header = map[string]interface{}{
		"typ": c.Opts.HeaderOpts.Type,
		"alg": c.Opts.HeaderOpts.Alg,
		"kid": c.Opts.HeaderOpts.KID,
	}
	claims := jwtgo.MapClaims{
		"iss": c.Opts.ClaimsOpts.Issuer,       // who creates the token and signs it
		"aud": c.Opts.ClaimsOpts.Audience,     // whom is allowed to use the token
		"exp": c.Opts.ClaimsOpts.Expiration(), // time when the token will expire (unix timestamp)
		"jti": c.Opts.ClaimsOpts.JWTID(),      // a unique identifier for the token
		"iat": c.Opts.ClaimsOpts.IssuedAt(),   // when the token was issued/created (now)
		"nbf": c.Opts.ClaimsOpts.NotBefore(),  // time before which the token is not yet valid (2 minutes ago)
	}
	for customClaimKey, customClaimVal := range customClaims {
		claims[customClaimKey] = customClaimVal
	}
	token.Claims = claims
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
