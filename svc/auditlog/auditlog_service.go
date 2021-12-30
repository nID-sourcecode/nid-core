package main

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/golang/protobuf/ptypes/empty"

	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	pb "lab.weave.nl/nid/nid-core/svc/auditlog/proto"
)

// AuditLogServiceServer implements the proto.AuditLogServiceServer interface
type AuditLogServiceServer struct {
	// Set the logger here so we can change it's output in the future to a file for example
	logger log.LoggerUtility
	stats  *Stats
}

// LogRequest implements the LogRequest grpc call, it simply logs the request message it receives
func (a *AuditLogServiceServer) LogRequest(ctx context.Context, request *pb.Request) (*empty.Empty, error) {
	var claims interface{}
	auth := request.GetAuth()
	if len(auth) > 0 {
		claims = decodeToken(request.Auth)
	}

	url, err := url.PathUnescape(request.GetUrl())
	if err != nil {
		log.WithField("url", url).WithError(err).Error("unable to unescape path")
		return nil, grpcerrors.ErrInvalidArgument("url not pathunescapable")
	}

	a.logger.WithFields(log.Fields{
		// Just print the claims as an object, let the log formatter handle the formatting
		"token":       claims,
		"url":         url,
		"body":        request.GetBody(),
		"http_method": request.GetHttpMethod(),
		"request_id":  request.GetRequestId(),
	}).Info("received response")

	return &empty.Empty{}, nil
}

// LogResponse auditlogs a response (designed to be called from the auditlog filter)
func (a *AuditLogServiceServer) LogResponse(ctx context.Context, response *pb.Response) (*empty.Empty, error) {
	a.logger.WithFields(log.Fields{
		"request_id":  response.RequestId,
		"status_code": response.StatusCode,
	}).Info("received request")

	return &empty.Empty{}, nil
}

func decodeToken(token string) interface{} {
	var claims interface{}
	rawTokenBody, err := base64url.Decode(token)
	if err != nil {
		claims = "token not parsable"
	} else {
		decodedToken := map[string]interface{}{}

		err = json.Unmarshal(rawTokenBody, &decodedToken)
		if err != nil {
			log.WithField("raw_token_body", rawTokenBody).Info("unable to json unmarshal token body token")
			// Do not return an error, just log that it is invalid
			claims = "token not parsable"
		} else {
			// Only print claims, we do not want to print the actual signature since this would be a security risk
			claims = decodedToken
		}
	}
	return claims
}
