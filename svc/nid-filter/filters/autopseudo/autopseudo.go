// Package autopseudo contains the autopseudo filter logic
package autopseudo

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/golang-jwt/jwt"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
)

type accessClaims struct {
	Subjects map[string]string
}

func (claims accessClaims) Valid() error {
	return nil
}

const (
	bearerScheme = "Bearer "
)

// Config contains the autopseudo filter configuratiion
type Config struct {
	FilterName        string
	Namespace         string
	TranslateToBSN    bool
	SubjectIdentifier string
	Key               *rsa.PrivateKey
	WalletClient      proto.WalletClient
}

// New creates a new filter initializer
func New(config *Config) *Filter {
	return &Filter{config: config}
}

// Filter is responsible for processing a single HTTP request and response
type Filter struct {
	config      *Config
	authHeader  *string
	replacement *string
}

var (
	errNoPathHeader              = errors.New("no :path header found")
	errNoAuthorizationHeader     = errors.New("no authorization header found")
	errNoReplacementFoundInToken = errors.New("no replacement found in the token, you may be accessing the wrong service")
)

// Check processes the HTTP request
func (a *Filter) Check(ctx context.Context, request *authv3.CheckRequest) error {
	headers := request.GetAttributes().GetRequest().GetHttp().GetHeaders()

	if authHeader, ok := headers["authorization"]; ok {
		a.authHeader = &authHeader
	}

	// handle headers
	path, ok := headers[":path"]
	if !ok {
		return errNoPathHeader
	}
	log.Infof("path: %s", path)
	queryIndex := strings.Index(path, "?")
	hasQuery := queryIndex != -1

	isIdentifierInQuery := false
	var query url.Values

	if hasQuery {
		queryString := path[(queryIndex + 1):]
		log.Debugf("querystring: %s", path)

		var err error
		query, err = url.ParseQuery(queryString)
		if err != nil {
			return errors.Wrap(err, "parsing query")
		}

		isIdentifierInQuery = a.isIdentifierInQuery(query)
	}

	body := request.GetAttributes().GetRequest().GetHttp().GetBody()

	if isIdentifierInQuery || body != "" {
		err := a.parseReplacement(ctx)
		if err != nil {
			log.Errorf("parse pseudo error or res, %+v", err)
			return err
		}
	}

	if isIdentifierInQuery {
		for _, v := range query {
			for i, val := range v {
				if strings.Contains(val, a.config.SubjectIdentifier) {
					v[i] = strings.ReplaceAll(val, a.config.SubjectIdentifier, *a.replacement)
				}
			}
		}
		adjustedQuery := query.Encode()
		adjustedPath := path[:(queryIndex+1)] + adjustedQuery
		headers[":path"] = adjustedPath
		log.Infof("adjustedPath: %s", adjustedPath)
	}

	var newBody string
	if body != "" {
		if strings.Contains(body, a.config.SubjectIdentifier) {
			newBodyString := strings.ReplaceAll(body, a.config.SubjectIdentifier, *a.replacement)
			headers["content-length"] = fmt.Sprintf("%d", len(newBodyString))
			newBody = newBodyString
		}
	}

	request.GetAttributes().GetRequest().GetHttp().Body = newBody
	request.GetAttributes().GetRequest().GetHttp().Headers = headers

	return nil
}

// Name returns the filter name
func (a *Filter) Name() string {
	return a.config.FilterName
}

// Perhaps turn this into a package?
func (a *Filter) decryptPseudo(encryptedPseudo string) (string, error) {
	encryptedPseudoBytes, err := base64.StdEncoding.DecodeString(encryptedPseudo)
	if err != nil {
		return "", fmt.Errorf("encrypted replacement base64 decode error: %w", err)
	}

	decryptedPseudoBytes, err := rsa.DecryptPKCS1v15(rand.Reader, a.config.Key, encryptedPseudoBytes)
	if err != nil {
		return "", fmt.Errorf("error decrypting replacement: %w", err)
	}

	decryptedPseudo := base64.StdEncoding.EncodeToString(decryptedPseudoBytes)

	return decryptedPseudo, nil
}

func (a *Filter) isIdentifierInQuery(query url.Values) bool {
	for _, v := range query {
		for _, val := range v {
			if strings.Contains(val, a.config.SubjectIdentifier) {
				return true
			}
		}
	}

	return false
}

func parseToken(authHeader string) (*accessClaims, error) {
	token := authHeader[len(bearerScheme):]
	claims := accessClaims{}

	jwtParser := jwt.Parser{}
	if _, _, err := jwtParser.ParseUnverified(token, &claims); err != nil {
		return nil, errors.Wrap(err, "unable to parse unverified JWT")
	}

	return &claims, nil
}

func (a *Filter) parseReplacement(ctx context.Context) error {
	if a.authHeader == nil {
		return errNoAuthorizationHeader
	}

	claims, err := parseToken(*a.authHeader)
	if err != nil {
		return errors.Wrap(err, "parsing token")
	}

	encryptedPseudonym, ok := claims.Subjects[a.config.Namespace]
	if !ok {
		return errNoReplacementFoundInToken
	}

	pseudonym, err := a.decryptPseudo(encryptedPseudonym)
	if err != nil {
		return errors.Wrap(err, "decrypting replacement")
	}

	if !a.config.TranslateToBSN {
		a.replacement = &pseudonym

		return nil
	}

	res, err := a.config.WalletClient.GetBSNForPseudonym(ctx, &proto.GetBSNForPseudonymRequest{Pseudonym: pseudonym})
	if err != nil {
		return errors.Wrap(err, "getting bsn for pseudonym from wallet")
	}
	bsn := res.GetBsn()
	a.replacement = &bsn
	return nil
}
