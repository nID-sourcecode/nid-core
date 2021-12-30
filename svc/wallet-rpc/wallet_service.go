package main

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/models"
	"lab.weave.nl/nid/nid-core/svc/wallet-rpc/gqlclient"
	"lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
)

// WalletServer handles authorization for the dashboard
type WalletServer struct {
	db    *WalletDB
	stats *Stats

	// Clients
	authClient gqlclient.IAuthClient
}

// ErrUserNotFound is returned when user was not found by pseudonym
var ErrUserNotFound = grpcerrors.ErrNotFound("user not found")

// CreateConsent will create a new client
func (w *WalletServer) CreateConsent(ctx context.Context, in *proto.CreateConsentRequest) (*proto.ConsentResponse, error) {
	var err error
	var consent *models.Consent
	var out *proto.ConsentResponse

	err = database.Transact(w.db.db, func(tx *gorm.DB) error {
		clientID := uuid.FromStringOrNil(in.ClientId)

		// Get userID from user pseudo
		user, err := models.NewUserDB(tx).GetByPseudo(in.UserPseudo)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			log.Extract(ctx).WithError(err).Error("unable to get user by pseudo")
			return ErrInternal
		}

		client, err := w.createOrFindClient(ctx, tx, clientID)
		if err != nil {
			log.Extract(ctx).WithError(err).WithField("client_id", clientID).Error("unable to find or create client")
			return ErrInternal
		}

		consent, err = w.createToConsentModel(in, client.ID, user.ID)
		if err != nil {
			log.Extract(ctx).WithError(err).Error("unable to create consent model")
			return ErrInternal
		}

		if err := models.NewConsentDB(tx).Add(ctx, consent); err != nil {
			log.Extract(ctx).WithError(err).Error("unable to store consent in db")
			return ErrInternal
		}

		// Create the ConsentResponse
		out, err = w.consentToConsentResponse(consent)
		if err != nil {
			log.Extract(ctx).WithError(err).Error("unable to create consent response")
			return ErrInternal
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"user_id":     consent.UserID,
		"client_id":   consent.ClientID,
		"description": consent.Description,
		"granted_at":  consent.Granted,
	}).Info("consent given")

	return out, nil
}

// GetBSNForPseudonym returns the bsn corresponding to a pseudonym
// FIXME this kind of thing should actually just be GraphQL. But what we need for that is static GQL golang client generation + easier istio auth on gql queries
func (w *WalletServer) GetBSNForPseudonym(ctx context.Context, req *proto.GetBSNForPseudonymRequest) (*proto.GetBSNForPseudonymResponse, error) {
	user, err := w.db.UserDB.GetByPseudo(req.GetPseudonym())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, grpcerrors.ErrNotFound("user not found")
		}
		log.Extract(ctx).WithError(err).Error("getting user by pseudonym")
		return nil, grpcerrors.ErrInternalServer()
	}

	return &proto.GetBSNForPseudonymResponse{
		Bsn: user.Bsn,
	}, nil
}
