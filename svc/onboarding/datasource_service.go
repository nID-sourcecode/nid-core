package main

import (
	context "context"
	"fmt"

	"lab.weave.nl/nid/nid-core/pkg/istioutil"
	"lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	onboardingPB "lab.weave.nl/nid/nid-core/svc/onboarding/proto"
	pseudoPB "lab.weave.nl/nid/nid-core/svc/pseudonymization/proto"
)

// DataSourceServiceServer server for handling onboarding of data sources
type DataSourceServiceServer struct {
	stats                  *Stats
	walletClient           gqlclient.Client
	pseudonimizationClient pseudoPB.PseudonymizerClient
	metadataHelper         headers.MetadataHelper
}

// Response Users response from wallet
type Response struct {
	Users []PseudoResponse `json:"users"`
}

// PseudoResponse response from wallet containing the pseudonym related to the BSN
type PseudoResponse struct {
	Pseudonym string `json:"pseudonym"`
}

// ConvertBSNToPseudonym converts a bsn to pseudonym for target namespace.
func (s *DataSourceServiceServer) ConvertBSNToPseudonym(ctx context.Context, in *onboardingPB.ConvertMessage) (*onboardingPB.ConvertResponseMessage, error) {
	certHeader, err := s.metadataHelper.GetValFromCtx(ctx, "x-forwarded-client-cert")
	if err != nil {
		return nil, grpcerrors.ErrInvalidArgument("missing header x-forwarded-client-cert")
	}

	fromNamespace, err := istioutil.GetNamespaceFromCertificateHeader(certHeader)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to get namespace from certificate header")

		return nil, grpcerrors.ErrInternalServer()
	}

	request := gqlclient.NewRequest(fmt.Sprintf(`{
		users(filter: {
			bsn: {
				eq: "%s"
			}
		}){
			pseudonym 
		}
	}
	  `, in.GetBsn()))
	res := Response{}

	err = s.walletClient.Post(ctx, request, &res)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to retrieve bsn from wallet")
		return nil, grpcerrors.ErrInternalServer()
	}
	if len(res.Users) != 1 {
		return nil, grpcerrors.ErrNotFound("bsn not found")
	}

	convertResponse, err := s.pseudonimizationClient.Convert(ctx, &pseudoPB.ConvertRequest{
		NamespaceTo: fromNamespace,
		Pseudonyms:  []string{res.Users[0].Pseudonym},
	})
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to convert pseudonym to target namespace")
		return nil, grpcerrors.ErrInternalServer()
	}

	result, ok := convertResponse.Conversions[res.Users[0].Pseudonym]
	if !ok {
		log.Extract(ctx).Error("pseudonym translation not found in conversion result")
		return nil, grpcerrors.ErrInternalServer()
	}

	return &onboardingPB.ConvertResponseMessage{
		Pseudonym: result,
	}, nil
}
