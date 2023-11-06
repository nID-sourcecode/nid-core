// Package autobsn
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	errgrpc "github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/autobsn/proto"
	walletPB "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
)

// AutoBSNServer the server for autobsn
type AutoBSNServer struct {
	Key          *rsa.PrivateKey
	walletClient walletPB.WalletClient
}

// NewAutoBSNServer creates a new autobsn server
func NewAutoBSNServer(key *rsa.PrivateKey, walletClient walletPB.WalletClient) *AutoBSNServer {
	return &AutoBSNServer{
		Key:          key,
		walletClient: walletClient,
	}
}

type accessClaims struct {
	Sub string
}

func (claims accessClaims) Valid() error {
	return nil
}

// ReplacePlaceholderWithBSN replaces the $$nid:bsn$$ placeholder in the request with the BSN that corresponds to the JWT
func (a *AutoBSNServer) ReplacePlaceholderWithBSN(ctx context.Context, req *proto.ReplacePlaceholderWithBSNRequest) (*proto.ReplacePlaceholderWithBSNResponse, error) {
	body := req.GetBody()

	query, err := url.ParseQuery(req.GetQuery())
	if err != nil {
		return nil, errgrpc.ErrInvalidArgument(errors.Wrap(err, "invalid query string").Error())
	}

	if !strings.Contains(body, subjectIdentifier) && !isIdentifierInQuery(query) {
		return &proto.ReplacePlaceholderWithBSNResponse{
			Body:  body,
			Query: req.GetQuery(),
		}, nil
	}

	authHeader := req.GetAuthorizationHeader()

	claims, err := parseToken(authHeader)
	if err != nil {
		return nil, errgrpc.ErrInvalidArgument(err.Error())
	}

	if claims.Sub == "" {
		return nil, errgrpc.ErrInvalidArgument("token does not contain subject")
	}

	encryptedPseudo := claims.Sub
	decryptedPseudo, err := a.decryptPseudo(encryptedPseudo)
	if err != nil {
		return nil, errgrpc.ErrInvalidArgument(err.Error())
	}

	resp, err := a.walletClient.GetBSNForPseudonym(ctx, &walletPB.GetBSNForPseudonymRequest{Pseudonym: decryptedPseudo})
	if err != nil {
		log.WithError(err).Error("getting bsn from wallet")
		return nil, errgrpc.ErrInternalServer()
	}
	bsn := resp.GetBsn()

	adjustedBody := strings.ReplaceAll(body, subjectIdentifier, bsn)

	for _, v := range query {
		for i, val := range v {
			if strings.Contains(val, subjectIdentifier) {
				v[i] = strings.ReplaceAll(val, subjectIdentifier, bsn)
			}
		}
	}
	adjustedQuery := query.Encode()

	return &proto.ReplacePlaceholderWithBSNResponse{
		Body:  adjustedBody,
		Query: adjustedQuery,
	}, nil
}

func (a *AutoBSNServer) decryptPseudo(encryptedPseudo string) (string, error) {
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
