package app

import (
	"encoding/json"

	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

func setClientMetadataToClaims(session *models.Session) (map[string]interface{}, error) {
	clientMetadataJSON := session.Client.Metadata.RawMessage
	clientMetadata := make(map[string]interface{})
	var err error
	if len(clientMetadataJSON) > 0 {
		err = json.Unmarshal(clientMetadataJSON, &clientMetadata)
	}
	return clientMetadata, err
}
