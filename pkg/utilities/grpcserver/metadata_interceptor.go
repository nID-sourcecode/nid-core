package grpcserver

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/grpckeys"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/logfields"
)

// unaryContextLogInterceptor adds a request id and a logger to the context
func unaryContextLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		ctx = ctxlogrus.ToContext(ctx, logrus.WithFields(logrus.Fields{
			logfields.MetadataError: "unable to get metadata from context",
			logfields.Context:       ctx,
		}))
	}

	ctx = addCtxTagsFromMetadata(ctx, meta)

	return handler(ctx, req)
}

// addCtxTagsFromMetadata adds tags to the request, if it has an error, the error is added to the logfields on the context
func addCtxTagsFromMetadata(ctx context.Context, metadata metadata.MD) context.Context {
	ctxTags := grpcctxtags.Extract(ctx)

	setRequestID(ctxTags, metadata)

	if envoyPath := metadata[grpckeys.EnvoyPathKey.String()]; len(envoyPath) > 0 && !ctxTags.Has(grpckeys.EnvoyPathKey.String()) {
		ctxTags.Set(grpckeys.EnvoyPathKey.String(), envoyPath[0])
	}

	if !ctxTags.Has(grpckeys.UserIDKey.String()) {
		// At this point, all jwt's should be valid (they are checked in istio), we only have to decode the claims
		claims, present, err := checkJWTHeader(metadata)

		if present && err != nil {
			// If we cannot get the claims, add that to the logfields
			ctx = ctxlogrus.ToContext(ctx, logrus.WithFields(logrus.Fields{
				logfields.BearerError: "unable to get metadata from context",
				logfields.Bearer:      metadata[grpckeys.AuthorizationKey.String()][0],
			}))
		}
		if present && err == nil {
			addClaimsToCtxTags(ctxTags, claims)
		} else {
			ctx = checkAccountMetadata(ctx, metadata, ctxTags)
		}
	}

	ctx = addCtxTagsToCtx(ctx, ctxTags)
	ctx = addIPAddressToCtx(ctx)

	return ctx
}

func addCtxTagsToCtx(ctx context.Context, ctxTags grpcctxtags.Tags) context.Context {
	// If the ctxTag is present, add the userID to the context
	if ctxTags.Has(grpckeys.UserIDKey.String()) {
		userID := ctxTags.Values()[grpckeys.UserIDKey.String()]
		ctx = context.WithValue(ctx, grpckeys.UserIDKey, userID)
		ctx = metadata.AppendToOutgoingContext(ctx, grpckeys.UserIDKey.String(), fmt.Sprintf("%v", userID))
	}

	// If the ctxTag is present, add the accountID to the context
	if ctxTags.Has(grpckeys.AccountIDKey.String()) {
		accountID := ctxTags.Values()[grpckeys.AccountIDKey.String()]
		ctx = context.WithValue(ctx, grpckeys.AccountIDKey, accountID)
		ctx = metadata.AppendToOutgoingContext(ctx, grpckeys.AccountIDKey.String(), fmt.Sprintf("%v", accountID))
	}

	return ctx
}

func addIPAddressToCtx(ctx context.Context) context.Context {
	metaHelper := &headers.GRPCMetadataHelper{}
	ip, err := metaHelper.GetIPFromCtx(ctx)
	if err != nil {
		return ctx
	}

	ctx = context.WithValue(ctx, grpckeys.IPAddress, ip)
	return metadata.AppendToOutgoingContext(ctx, grpckeys.IPAddress.String(), fmt.Sprintf("%v", ip))
}

func setRequestID(ctxTags grpcctxtags.Tags, metadata metadata.MD) {
	if !ctxTags.Has(grpckeys.DefaultXRequestIDKey.String()) {
		if requestID := metadata[grpckeys.RequestIDKey.String()]; len(requestID) > 0 {
			ctxTags.Set(grpckeys.DefaultXRequestIDKey.String(), requestID[0])
		} else {
			// Add requestID to the tags
			ctxTags.Set(grpckeys.DefaultXRequestIDKey.String(), shortID())
		}
	}
}

func addClaimsToCtxTags(ctxTags grpcctxtags.Tags, claims jwt.MapClaims) {
	// Add the userID to the tags and context
	ctxTags.Set(grpckeys.UserIDKey.String(), claims["sub"])

	// Add the accountID to the tags and context
	ctxTags.Set(grpckeys.AccountIDKey.String(), claims["account_id"])

	// Add the apiKeyID to the tags if it is present
	apiKeyID := claims["api_key_id"]
	if apiKeyID != nil {
		ctxTags.Set(grpckeys.APIKeyIDKey.String(), apiKeyID)
	}
}

func checkAccountMetadata(ctx context.Context, metadata metadata.MD, ctxTags grpcctxtags.Tags) context.Context {
	accountID := metadata.Get(grpckeys.AccountIDKey.String())
	if len(accountID) == 1 {
		ctxTags.Set(grpckeys.AccountIDKey.String(), accountID[0])
		// At this point we reached a route which no JWT was set, but RBAC
		// allowed the request which means this must be an internal server call
		// or a public route. Now if the account id is set, use it on the context
		ctx = context.WithValue(ctx, grpckeys.AccountIDKey, ctxTags.Values()[grpckeys.AccountIDKey.String()])
	}

	userID := metadata.Get(grpckeys.UserIDKey.String())
	if len(userID) == 1 {
		ctxTags.Set(grpckeys.UserIDKey.String(), userID[0])
		ctx = context.WithValue(ctx, grpckeys.UserIDKey, ctxTags.Values()[grpckeys.UserIDKey.String()])
	}

	return ctx
}

func checkJWTHeader(metadata metadata.MD) (jwt.MapClaims, bool, error) {
	if bearer := metadata[grpckeys.AuthorizationKey.String()]; len(bearer) > 0 && strings.HasPrefix(bearer[0], BearerPrefix) {
		claims, err := GetClaimsWithoutValidation(strings.TrimPrefix(bearer[0], BearerPrefix))
		if err != nil {
			return nil, true, err
		}

		return claims, true, nil
	}

	return nil, false, nil
}

// GetUserIDFromContext returns the userID if it is present on the context.
// If the userID is not the correct format or not parsable, the function panics
// due to the assumption that the userID should be present or correct when performing this action in a route
func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	var rawUserID string
	var ok bool
	if rawUserID, ok = ctx.Value(grpckeys.UserIDKey).(string); !ok {
		panic("user_id should be present on the context and be of type string")
	}

	userID, err := uuid.FromString(rawUserID)
	if err != nil {
		panic(fmt.Errorf("unable to parse user_id from context: %w", err))
	}

	return userID
}

// GetAccountIDFromContext returns the accountID if it is present on the context.
// If the accountID is not the correct format or not parsable, the function panics
// due to the assumption that the accountID should be present or correct when performing this action in a route
func GetAccountIDFromContext(ctx context.Context) string {
	var rawAccountID string
	var ok bool
	if rawAccountID, ok = ctx.Value(grpckeys.AccountIDKey).(string); !ok {
		panic("account_id should be present on the context and be of type string")
	}
	return rawAccountID
}

// GetIPAddressFromContext returns the IP address if it is present on the context.
// If the ipAddress is not the correct format or not parsable, the function returns a empty string
func GetIPAddressFromContext(ctx context.Context) string {
	var address string
	var ok bool
	if address, ok = ctx.Value(grpckeys.IPAddress).(string); !ok {
		return ""
	}
	return address
}

// shortID generates a random id of length 6
func shortID() string {
	b := make([]byte, 6)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic("unable to read bytes for shortID")
	}
	return base64.StdEncoding.EncodeToString(b)
}
