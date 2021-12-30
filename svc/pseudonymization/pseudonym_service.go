package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"

	"lab.weave.nl/nid/nid-core/pkg/istioutil"
	goErr "lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/pseudonymization/keymanager"
	pb "lab.weave.nl/nid/nid-core/svc/pseudonymization/proto"
	"lab.weave.nl/nid/nid-core/svc/pseudonymization/pseudonymization"
)

// PseudonymizerServer manages all the routes for Pseudonymization
type PseudonymizerServer struct {
	stats          *Stats
	KeyManager     keymanager.KeyManager
	metadataHelper headers.MetadataHelper
}

// Generate generates given amount of pseudonyms
func (p PseudonymizerServer) Generate(ctx context.Context, req *pb.GenerateRequest) (*pb.GenerateResponse, error) {
	pseudonyms := make([]string, req.GetAmount())
	var i uint32
	for i = 0; i < req.Amount; i++ {
		pseudonym, err := generate()
		if err != nil {
			log.Extract(ctx).WithError(err).WithField("amount", req.GetAmount()).Error("unable to generate pseudonym")

			return nil, errors.ErrInternalServer()
		}
		pseudonyms[i] = pseudonym
	}

	return &pb.GenerateResponse{Pseudonyms: pseudonyms}, nil
}

func generate() (string, error) {
	internalID := make([]byte, 32)
	if _, err := rand.Read(internalID); err != nil {
		return "", goErr.Wrap(err, "unable to generate random pseudonym")
	}

	pseudonymBytes, err := pseudonymization.Encode(internalID, make([]byte, 32))
	if err != nil {
		return "", goErr.Wrap(err, "unable to encode pseudonymization")
	}

	return base64.StdEncoding.EncodeToString(pseudonymBytes), nil // Can use empty service name since it is irrelevant at this stage
}

// Convert converts pseudonyms to given namespace
func (p PseudonymizerServer) Convert(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	certHeader, err := p.metadataHelper.GetValFromCtx(ctx, "x-forwarded-client-cert")
	if err != nil {
		return nil, errors.ErrInvalidArgument("missing header x-forwarded-client-cert")
	}

	fromNamespace, err := istioutil.GetNamespaceFromCertificateHeader(certHeader)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to get namespace from certificate header")

		return nil, errors.ErrInternalServer()
	}

	rsaKey, err := p.KeyManager.GetKey(req.GetNamespaceTo())
	if err != nil {
		log.Extract(ctx).WithError(err).WithField("namespace", req.GetNamespaceTo()).Error("unable to get key for namespace")

		return nil, errors.ErrInternalServer()
	}

	// Convert pseudonyms
	conversions := make(map[string][]byte)
	for _, pseudonymIn := range req.GetPseudonyms() {
		pseudonymOut, err := convert(pseudonymIn, fromNamespace, req.GetNamespaceTo())
		if err != nil {
			log.Extract(ctx).WithError(err).WithFields(log.Fields{
				"pseudonyms":     pseudonymIn,
				"namespace_to":   req.GetNamespaceTo(),
				"namespace_from": fromNamespace,
			}).Error("unable to get key for namespace")

			return nil, errors.ErrInternalServer()
		}
		encryptedPseudonymOut, err := rsa.EncryptPKCS1v15(rand.Reader, rsaKey, pseudonymOut)
		if err != nil {
			log.Extract(ctx).WithError(err).WithField("pseudonym", pseudonymOut).Error()

			return nil, errors.ErrInternalServer()
		}
		conversions[pseudonymIn] = encryptedPseudonymOut
	}

	return &pb.ConvertResponse{Conversions: conversions}, nil
}

func convert(pseudonym, serviceFromName, serviceToName string) ([]byte, error) {
	internalID, err := pseudonymization.Decode(pseudonym, []byte(serviceFromName))
	if err != nil {
		return nil, goErr.Wrap(err, "unable to decode pseudonym from name")
	}

	return pseudonymization.Encode(internalID, []byte(serviceToName))
}
