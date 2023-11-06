package http

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
)

// GetBasicAuth retrieve basic auth from context
func getBasicAuth(c *gin.Context) (string, string, error) {
	basicAuth := c.GetHeader("Authorization")
	if basicAuth == "" {
		return "", "", errors.Wrapf(contract.ErrUnauthenticated, "getting authorization header from request")
	}

	return parseBasicAuth(basicAuth)
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok error) {
	const prefix = "Basic "
	// Case-insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", "", contract.ErrIncorrectBasicAuthPrefix
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", fmt.Errorf("error decoding basic auth: %w", err)
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return "", "", contract.ErrIncorrectBasicAuthFormat
	}
	return cs[:s], cs[s+1:], nil
}
