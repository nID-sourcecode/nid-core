package graphql

import (
	"lab.weave.nl/nid/nid-core/svc/info-manager/proto"
)

// InfoManagerClient inits the proto client.
var InfoManagerClient proto.InfoManagerClient //nolint:gochecknoglobals

// Init takes and sets the client.
func Init(client proto.InfoManagerClient) {
	InfoManagerClient = client
}
