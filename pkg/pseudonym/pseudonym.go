// Package pseudonym provides functionality for retrieving and generating pseudonyms
package pseudonym

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	pb "lab.weave.nl/nid/nid-core/svc/pseudonymization/proto"
)

// Only for non-istio debug purposes
const (
	certMetatadataKey       = "x-forwarded-client-cert"
	certFormat              = "SAN=\"\";URI=spiffe://cluster.local/ns/%s/sa/testservice"
	defaultTestingNamespace = "alice"
)

// Pseudonymizer translates and generates pseudonyms
type Pseudonymizer interface {
	GetPseudonym(ctx context.Context, myPseudo, targetNamespace string) (string, error)
	GeneratePseudonym(ctx context.Context, amount uint32) ([]string, error)
}

type pseudonymizer struct {
	location         string
	istioDisabled    bool
	testingNamespace string
}

// NewPseudonymizer creates a new pseudonymizer
func NewPseudonymizer(location string) Pseudonymizer {
	istioDisabled := os.Getenv("NO_ISTIO") == "TRUE"
	testingNamespace := os.Getenv("TESTING_NAMESPACE")
	if testingNamespace == "" {
		testingNamespace = defaultTestingNamespace
	}

	return &pseudonymizer{location: location, istioDisabled: istioDisabled, testingNamespace: testingNamespace}
}

func (w pseudonymizer) GetPseudonym(ctx context.Context, myPseudo, targetNamespace string) (string, error) {
	conn, err := grpc.Dial(w.location, grpc.WithInsecure())
	if err != nil {
		return "", errors.Wrap(err, "unable to dial pseudonym service")
	}
	defer func() {
		log.Error(conn.Close())
	}()
	c := pb.NewPseudonymizerClient(conn)

	req := pb.ConvertRequest{Pseudonyms: []string{myPseudo}, NamespaceTo: targetNamespace}

	if w.istioDisabled { // Fake client certificate
		md := metadata.Pairs(certMetatadataKey, fmt.Sprintf(certFormat, w.testingNamespace))
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	res, err := c.Convert(ctx, &req)
	if err != nil {
		return "", errors.Wrap(err, "unable to convert pseudonym")
	}

	pseudoString := base64.StdEncoding.EncodeToString(res.Conversions[myPseudo])

	return pseudoString, nil
}

func (w pseudonymizer) GeneratePseudonym(ctx context.Context, amount uint32) ([]string, error) {
	conn, err := grpc.Dial(w.location, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "unable to dial pseudonym service")
	}
	defer func() {
		log.Error(conn.Close())
	}()
	c := pb.NewPseudonymizerClient(conn)

	req := pb.GenerateRequest{Amount: amount}

	if w.istioDisabled { // Fake client certificate
		md := metadata.Pairs(certMetatadataKey, fmt.Sprintf(certFormat, w.testingNamespace))
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	res, err := c.Generate(ctx, &req)
	if err != nil {
		return nil, errors.Wrap(err, "pseudo service was unable to generate pseudonym")
	}

	return res.Pseudonyms, nil
}
