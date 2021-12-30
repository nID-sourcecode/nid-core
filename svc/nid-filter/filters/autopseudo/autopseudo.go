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

	"github.com/dgrijalva/jwt-go"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"

	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/nid-filter/filters/utils"
	"lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
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

// FilterInitializer creates new filters
type FilterInitializer struct {
	config *Config
}

// Name returns the filter name
func (a *FilterInitializer) Name() string {
	return a.config.FilterName
}

// NewFilterInitializer creates a new filter initializer
func NewFilterInitializer(config *Config) *FilterInitializer {
	return &FilterInitializer{config: config}
}

// NewFilter creates a new filter
func (a *FilterInitializer) NewFilter() (filter.Filter, error) {
	return &Filter{config: a.config}, nil
}

// Filter is responsible for processing a single HTTP request and response
type Filter struct {
	filter.DefaultFilter
	config      *Config
	authHeader  *string
	replacement *string
}

var errNoPathHeader = errors.New("no :path header found")

// OnHTTPRequest processes an HTTP request
func (a *Filter) OnHTTPRequest(ctx context.Context, body []byte, headers map[string]string) (*filter.ProcessingResponse, error) {
	if authHeader, ok := headers["authorization"]; ok {
		a.authHeader = &authHeader
	}

	// handle headers
	path, ok := headers[":path"]
	if !ok {
		return nil, errNoPathHeader
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
			return nil, errors.Wrap(err, "parsing query")
		}

		isIdentifierInQuery = a.isIdentifierInQuery(query)
	}

	if isIdentifierInQuery || body != nil {
		res, err := a.parseReplacement(ctx)
		if res != nil || err != nil {
			log.Infof("parse pseudo error or res, %+v, %+v", res, err)
			return res, err
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

	var newBody []byte
	if body != nil {
		bodyString := string(body)
		if strings.Contains(bodyString, a.config.SubjectIdentifier) {
			if a.replacement == nil {
			}

			newBodyString := strings.ReplaceAll(bodyString, a.config.SubjectIdentifier, *a.replacement)
			newBody = []byte(newBodyString)
			headers["content-length"] = fmt.Sprintf("%d", len(newBody))
		}
	}

	return &filter.ProcessingResponse{
		NewHeaders: headers,
		NewBody:    newBody,
	}, nil
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

func (a *Filter) parseReplacement(ctx context.Context) (*filter.ProcessingResponse, error) {
	if a.authHeader == nil {
		return utils.GraphqlError("no authorization header specified", envoy_type_v3.StatusCode_BadRequest), nil
	}

	claims, err := parseToken(*a.authHeader)
	if err != nil {
		return utils.GraphqlError(errors.Wrap(err, "parsing token").Error(), envoy_type_v3.StatusCode_BadRequest), nil
	}

	encryptedPseudonym, ok := claims.Subjects[a.config.Namespace]
	if !ok {
		return utils.GraphqlError("no replacement found in the token, you may be accessing the wrong service", envoy_type_v3.StatusCode_BadRequest), nil
	}

	pseudonym, err := a.decryptPseudo(encryptedPseudonym)
	if err != nil {
		return nil, errors.Wrap(err, "decrypting replacement")
	}

	if !a.config.TranslateToBSN {
		a.replacement = &pseudonym

		return nil, nil
	}

	res, err := a.config.WalletClient.GetBSNForPseudonym(ctx, &proto.GetBSNForPseudonymRequest{Pseudonym: pseudonym})
	if err != nil {
		return nil, errors.Wrap(err, "getting bsn for pseudonym from wallet")
	}
	bsn := res.GetBsn()
	a.replacement = &bsn
	return nil, nil
}
