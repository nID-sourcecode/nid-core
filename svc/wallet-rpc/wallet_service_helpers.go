package main

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
	"github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
)

func (w *WalletServer) createOrFindClient(ctx context.Context, tx *gorm.DB, clientID uuid.UUID) (*models.Client, error) {
	clientDB := models.NewClientDB(tx)

	// Find or create client
	client, err := clientDB.GetByExtClientID(clientID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Extract(ctx).WithError(err).WithField("ext_client_id", clientID).Error("unable to get client by ext_client_id")
		return nil, ErrInternal
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Fetch the client from the auth service
		authClient, err := w.authClient.FetchClient(ctx, clientID)
		if err != nil {
			log.Extract(ctx).WithError(err).WithField("ext_client_id", clientID).Error("unable to fetch auth client")
			return nil, ErrInternal
		}

		client := &models.Client{
			Color:       authClient.Color,
			ExtClientID: clientID.String(),
			Icon:        authClient.Icon,
			Logo:        authClient.Logo,
			Name:        authClient.Name,
		}

		if err := clientDB.Add(ctx, client); err != nil {
			log.Extract(ctx).WithError(err).WithField("ext_client_id", clientID).Error("unable to create client")
			return nil, ErrInternal
		}
		return client, nil
	}

	return client, nil
}

func (w *WalletServer) createToConsentModel(in *proto.CreateConsentRequest, clientID, userID uuid.UUID) (*models.Consent, error) {
	consent := &models.Consent{
		AccessToken: in.AccessToken,
		Description: in.Description,
		ClientID:    clientID,
		UserID:      userID,
		Name:        in.Name,
	}

	if in.GrantedAt != nil {
		err := in.GrantedAt.CheckValid()
		if err != nil {
			return nil, errors.Wrap(err, "unable to convert proto timestamp")
		}
		granted := in.GrantedAt.AsTime()

		consent.Granted = &granted
	}

	return consent, nil
}

func (w *WalletServer) consentToConsentResponse(consent *models.Consent) (*proto.ConsentResponse, error) {
	out := &proto.ConsentResponse{
		Id:          consent.ID.String(),
		AccessToken: consent.AccessToken,
		ClientId:    consent.ClientID.String(),
		Description: consent.Description,
		Name:        consent.Name,
		UserId:      consent.UserID.String(),
	}

	if consent.Granted != nil {
		pGranted := timestamppb.New(*consent.Granted)
		out.GrantedAt = pGranted
		if !pGranted.IsValid() {
			return nil, errors.New("unable to convert proto timestamp")
		}
	}

	return out, nil
}
