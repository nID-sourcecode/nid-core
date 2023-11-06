// Package istioutil provides utility functionality for istio
package istioutil

import (
	"fmt"
	"strings"
)

var (
	// ErrNoURIFound no URI entry
	ErrNoURIFound = fmt.Errorf("no URI entry with namespace found in certificate")
)

// GetNamespaceFromCertificateHeader converts an istio certificate header to a namespace
func GetNamespaceFromCertificateHeader(certHeader string) (string, error) {
	certPairs := strings.Split(certHeader, ";")
	for _, certPair := range certPairs {
		if strings.HasPrefix(certPair, "URI=spiffe://cluster.local/ns/") {
			postNsPart := certPair[30:]

			return strings.Split(postNsPart, "/")[0], nil
		}
	}

	return "", ErrNoURIFound
}
