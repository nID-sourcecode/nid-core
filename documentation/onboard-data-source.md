# Onboarding data source

## Prerequisites
- A running service in your platform namespace

## Retrieving pseudonyms

When an existing data source connects to the cluster, pseudonyms for all BSN's should be retrieved. This is necessary, since the platform expects that a data source keeps track of pseudonyms for all of its data.

This can be done by first contacting `ConvertBSNToPseudonym` rpc on the onboarding service. The onboarding service is available at `onboarding.twi.svc.cluster.local`. To obtain a pseudonym, a BSN should be set in the `ConvertMessage`. Note that this service can only be contacted when  from within the cluster. The onboarding service will return a pseudonym in the `ConvertResponseMessage`, which is a byte array containing a RSA encrypted pseudonym. Afterwards, we can base64 encode the RSA encrypted pseudonym and let the autopseudo service in our namespace decrypt it. The autopseudo service is available at `autopseudo.<my_namespace>.svc.cluster.local` and can be contacted as `GET` REST call at `autopseudo.<my_namespace>.svc.cluster.local/decrypt`. The base64 encoded pseudonym is expected as query parameter in `pseudonym`. An example implementation in Go can be found [here](#example-request). 

Data initially related to the BSN should now be related to the newly obtained pseudonym. This way a data source is able to process a request for a given pseudonym instead of a BSN. 

For example, if we have a User table consisting of an `id` and a `bsn`, we should add a new column, see the example below.

#### Old User table
| id                                   | bsn       |
|--------------------------------------|-----------|
| 10ba3c1a-96e7-4d3f-9440-9cd42909d1fc | 123456789 |

#### New User table
| id                                   | bsn       | pseudonym                                                        |
|--------------------------------------|-----------|------------------------------------------------------------------|
| 10ba3c1a-96e7-4d3f-9440-9cd42909d1fc | 123456789 | z0uYcrd6FsEpYBkC9ujn5dDYK9iGDjrQnfQ8fJShjM+J3zuiAi6JdkTLkEuIshml |

### Proto specification
```proto
syntax = "proto3";
package onboarding;

import "google/api/annotations.proto";
import "google/protobuf/empty.proto";
import "lab.weave.nl/devops/proto-istio-auth-generator/proto/scope.proto";

option go_package = ".;proto";

service DataSourceService {

  // ConvertBSNToPseudonym
  //
  // ConvertBSNToPseudonym converts a bsn to pseudonym for target namespace.
  rpc ConvertBSNToPseudonym(ConvertMessage) returns (ConvertResponseMessage) {
    option (scopes.scope) = "convertbsn";
    option (google.api.http) = {
      get : "/v1/onboarding/datasource/convertbsntopseudonym"
    };
  }
}

message ConvertMessage {
  string bsn = 1;
}

message ConvertResponseMessage {
  bytes pseudonym = 1;
}
```
  
## Example request

An example request in Go is shown below.

```GOLANG
package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"<my project path to proto package>/proto"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial("onboarding.twi.cluster.svc.local", grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		// Handle error
		panic(err)
	}
	dataSourceClient := proto.NewDataSourceServiceClient(conn)
	res, err := dataSourceClient.ConvertBSNToPseudonym(ctx, &proto.ConvertMessage{
		// Dummy BSN
		Bsn: "123456789",
	})
	if err != nil {
		// Handle error
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "autopseudo.<my_namespace>.svc.cluster.local", nil)
	if err != nil {
		// Handle error
		panic(err)
	}
	req.URL.Query().Add("pseudonym", base64.StdEncoding.EncodeToString(res.Pseudonym))

	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		// Handle error
		panic(err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// Handle error
			fmt.Println("Unable to close response body")
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Resulting pseudonym: %s\n", body)

	// Add decrypted pseudonym to my favorite DB
}
```
