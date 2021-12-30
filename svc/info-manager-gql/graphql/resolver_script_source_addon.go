package graphql

import (
	"context"
	"io"
	"net/http"

	"github.com/jinzhu/gorm"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/info-manager/models"
	"lab.weave.nl/nid/nid-core/svc/info-manager/proto"
)

// AfterReadGetSignedURL ask info-manager grpc service for a signed url and return with model
func (h *CustomScriptSourceHooks) AfterReadGetSignedURL(ctx context.Context, tx *gorm.DB, model *models.ScriptSource) error {
	var err error
	rpc, err := InfoManagerClient.ScriptsGet(ctx, &proto.ScriptsGetRequest{
		ScriptId: model.ScriptID.String(),
		Version:  model.Version,
	})
	if err != nil {
		log.WithFields(log.Fields{"script_id": model.ScriptID, "version": model.Version}).WithError(err).Error("calling info manager client to fetch signed url for source")
		return errors.Wrap(err, "fetching signed script url from info manager client failed")
	}

	model.SignedURL = &rpc.SignedUrl

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rpc.SignedUrl, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"script_id": model.ScriptID, "version": model.Version}).WithError(err).Error("calling signed url to fetch raw script failed")
		return errors.Wrap(err, "fetching raw script from signed url failed")
	}
	defer func() {
		closeError := resp.Body.Close()
		if err != nil {
			err = errors.CombineErrors(errors.Wrap(closeError, "closing request body failed"), err)
		}
	}()
	body, err := io.ReadAll(resp.Body)

	bodyString := string(body)
	model.RawScript = &bodyString

	return nil
}
