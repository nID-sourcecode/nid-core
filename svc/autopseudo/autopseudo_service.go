package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
)

// AutoPseudoServer the server for autopseudo
type AutoPseudoServer struct {
	Key  *rsa.PrivateKey
	Conf *AutoPseudoConfig
}

// NewAutoPseudoServer creates a new autopseudo server
func NewAutoPseudoServer(key *rsa.PrivateKey, conf *AutoPseudoConfig) *AutoPseudoServer {
	return &AutoPseudoServer{
		Key:  key,
		Conf: conf,
	}
}

// DecryptAndApply decryptes pseudonym and applies it on request body or parameters
func (a *AutoPseudoServer) DecryptAndApply(c *gin.Context) {
	body := ""
	if c.Request.Method != "GET" {
		bodyBytes, err := c.GetRawData()
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, graphqlError(err.Error()))

			return
		}
		body = string(bodyBytes)
	}

	query := c.Request.URL.Query()

	defaultResponse := decryptAndApplyResponse{
		Body:  body,
		Query: query.Encode(),
	}

	if !strings.Contains(body, subjectIdentifier) && !isIdentifierInQuery(query) {
		c.PureJSON(http.StatusOK, defaultResponse)

		return
	}

	authHeader := c.GetHeader("Authorization")

	claims, err := parseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, graphqlError(err.Error()))

		return
	}

	encryptedPseudo := claims.Subjects[a.Conf.Namespace]
	decryptedPseudo, err := a.decryptPseudo(encryptedPseudo)
	if err != nil {
		c.JSON(http.StatusBadRequest, graphqlError(err.Error()))

		return
	}

	adjustedBody := strings.ReplaceAll(body, subjectIdentifier, decryptedPseudo)

	for _, v := range query {
		for i, val := range v {
			if strings.Contains(val, subjectIdentifier) {
				v[i] = strings.ReplaceAll(val, subjectIdentifier, decryptedPseudo)
			}
		}
	}
	adjustedQuery := query.Encode()

	c.PureJSON(http.StatusOK, decryptAndApplyResponse{
		Body:  adjustedBody,
		Query: adjustedQuery,
	})
}

// Decrypt decrypts given pseudonym
func (a *AutoPseudoServer) Decrypt(c *gin.Context) {
	keys, ok := c.Request.URL.Query()["pseudonym"]
	if !ok || len(keys) != 1 {
		c.JSON(http.StatusBadRequest, graphqlError("did not find query parameter pseudonym"))

		return
	}

	decryptedPseudo, err := a.decryptPseudo(keys[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, graphqlError(err.Error()))

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"decrypted_pseudonym": decryptedPseudo,
	})
}

// Perhaps turn this into a package?
func (a *AutoPseudoServer) decryptPseudo(encryptedPseudo string) (string, error) {
	encryptedPseudoBytes, err := base64.StdEncoding.DecodeString(encryptedPseudo)
	if err != nil {
		return "", fmt.Errorf("encrypted pseudonym base64 decode error: %w", err)
	}

	decryptedPseudoBytes, err := rsa.DecryptPKCS1v15(rand.Reader, a.Key, encryptedPseudoBytes)
	if err != nil {
		return "", fmt.Errorf("error decrypting pseudonym: %w", err)
	}

	decryptedPseudo := base64.StdEncoding.EncodeToString(decryptedPseudoBytes)

	return decryptedPseudo, nil
}

func isIdentifierInQuery(query url.Values) bool {
	for _, v := range query {
		for _, val := range v {
			if strings.Contains(val, subjectIdentifier) {
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

func graphqlError(message string) interface{} {
	return gin.H{
		"errors": []gin.H{
			{
				"message": message,
			},
		},
	}
}
