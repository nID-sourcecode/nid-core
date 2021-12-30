package headers

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/grpckeys"
)

// Error definitions
var (
	ErrIncorrectBasicAuthPrefix error = fmt.Errorf("incorrect basic auth prefix")
	ErrIncorrectBasicAuthFormat error = fmt.Errorf("incorrect basic auth format")
)

// GetBasicAuth retrieve basic auth from context
func (m GRPCMetadataHelper) GetBasicAuth(ctx context.Context) (string, string, error) {
	basicAuth, err := m.GetValFromCtx(ctx, grpckeys.AuthorizationKey.String())
	if err != nil {
		return "", "", err
	}
	return m.parseBasicAuth(basicAuth)
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func (m GRPCMetadataHelper) parseBasicAuth(auth string) (username, password string, ok error) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", "", ErrIncorrectBasicAuthPrefix
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", fmt.Errorf("error decoding basic auth: %w", err)
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return "", "", ErrIncorrectBasicAuthFormat
	}
	return cs[:s], cs[s+1:], nil
}
