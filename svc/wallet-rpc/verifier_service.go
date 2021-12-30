package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/messagebird"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/postmark"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/variables"
	"lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
)

const (
	errGetSession     = "error getting %s verification session for %s from database"
	errTokenMisMatch  = "token mismatch"
	errVerifier       = "error verifying"
	errDatabaseUpdate = "error setting verification for %s %s, error %s"
	errAlreadyValid   = "%s already verified"
	errWait           = "wait: %ds"
	errCreateToken    = "error creating verification_token for %s %s, error %s"
)

// errros
var (
	errTokenExpired      error = grpcerrors.ErrNotFound("token exipred")
	errRetryAlreadyValid error = grpcerrors.ErrFailedPrecondition(errAlreadyValid)
)

// VerifierServer handles authorization for the dashboard
type VerifierServer struct {
	stats *Stats
	db    *WalletDB

	emailVerifier postmark.EmailVerifier
	phoneVerifier messagebird.PhoneVerifier
}

// VerifyEmail with external library
func (e *VerifierServer) VerifyEmail(ctx context.Context, in *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	emailAddress, err := e.db.EmailAddressDB.Get(ctx, uuid.FromStringOrNil(in.GetId()))
	if err != nil {
		log.Extract(ctx).WithError(err).Errorf(errGetSession, "email address", in.GetId())
		return nil, ErrInternal
	}

	response := &proto.VerifyResponse{Id: in.GetId()}

	// Break if already verified
	if emailAddress.Verified {
		return response, nil
	}

	// Error if too much time has passed
	if emailAddress.UpdatedAt.Add(time.Second * variables.VerifyEmailAddressTimeout).Before(time.Now()) {
		return nil, errTokenExpired
	}

	// Verify code
	if err := e.emailVerifier.CheckEmailVerification(emailAddress.VerificationToken, in.GetCode()); err != nil {
		if errors.Is(err, postmark.ErrTokenDidNotMatch) {
			log.WithError(err).WithField("email_address_id", emailAddress.ID.String()).Errorf(errTokenMisMatch)
		} else {
			log.WithError(err).WithField("email_address_id", emailAddress.ID.String()).Errorf(errVerifier)
		}

		return nil, ErrInternal
	}

	// Save verification success to DB
	emailAddress.Verified = true
	emailAddress.UpdatedAt = time.Now()

	if err = e.db.EmailAddressDB.Update(ctx, emailAddress); err != nil {
		log.Extract(ctx).WithError(err).Errorf(errDatabaseUpdate, "email", emailAddress.ID.String(), err)
		return nil, ErrInternal
	}

	return response, nil
}

// RetryVerifyEmail restarts the verification process
func (e *VerifierServer) RetryVerifyEmail(ctx context.Context, in *proto.RetryVerifyRequest) (*proto.VerifyResponse, error) {
	emailAddress, err := e.db.EmailAddressDB.Get(ctx, uuid.FromStringOrNil(in.GetId()))
	if err != nil {
		log.Extract(ctx).WithError(err).Errorf(errGetSession, in.GetId())
		return nil, ErrInternal
	}

	// Check if email address needs to be verified
	if emailAddress.Verified {
		log.Extract(ctx).WithError(err).Errorf(errAlreadyValid, "email", emailAddress.ID.String(), err)
		return nil, errRetryAlreadyValid
	}

	// Check if enough time has passed to retry
	unitl := time.Until(emailAddress.UpdatedAt.Add(time.Second * variables.VerifyEmailAddressRetryDebounce))
	if dur := unitl; dur > 0 {
		return nil, grpcerrors.ErrFailedPrecondition(fmt.Sprintf(errWait, dur/time.Second))
	}

	// Get a new verification session
	token, err := e.emailVerifier.NewEmailVerification(emailAddress.EmailAddress)
	if err != nil {
		log.Extract(ctx).WithError(err).Errorf(errCreateToken, "email", emailAddress.ID.String(), err)
		return nil, ErrInternal
	}

	// Write verification session token
	emailAddress.UpdatedAt = time.Now()
	emailAddress.VerificationToken = token

	if err = e.db.EmailAddressDB.Update(ctx, emailAddress); err != nil {
		log.Extract(ctx).WithError(err).Errorf(errDatabaseUpdate, "email", emailAddress.ID.String(), err)
		return nil, ErrInternal
	}

	return &proto.VerifyResponse{Id: in.GetId()}, nil
}

// VerifyPhoneNumber verify phone number with external library
func (e *VerifierServer) VerifyPhoneNumber(ctx context.Context, in *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	phoneNumber, err := e.db.PhoneNumberDB.Get(ctx, uuid.FromStringOrNil(in.GetId()))
	if err != nil {
		log.Extract(ctx).WithError(err).Errorf(errGetSession, "phone number", in.GetId())
		return nil, ErrInternal
	}

	response := &proto.VerifyResponse{Id: in.GetId()}

	// Break if already verified
	if phoneNumber.Verified {
		return response, nil
	}

	// Verify code
	if err := e.phoneVerifier.CheckPhoneVerification(phoneNumber.VerificationToken, in.Code); err != nil {
		log.WithError(err).WithField("phone_number", phoneNumber.ID.String()).Errorf(errVerifier)

		return nil, ErrInternal
	}

	// Save verification success to DB
	phoneNumber.Verified = true
	phoneNumber.UpdatedAt = time.Now()

	if err = e.db.PhoneNumberDB.Update(ctx, phoneNumber); err != nil {
		log.Extract(ctx).WithError(err).Errorf(errDatabaseUpdate, "phone", phoneNumber.ID.String(), err)
		return nil, ErrInternal
	}

	return response, nil
}

// RetryVerifyPhoneNumber restarts the verification process
func (e *VerifierServer) RetryVerifyPhoneNumber(ctx context.Context, in *proto.RetryPhoneRequest) (*proto.VerifyResponse, error) {
	phoneNumber, err := e.db.PhoneNumberDB.Get(ctx, uuid.FromStringOrNil(in.GetId()))
	if err != nil {
		log.Extract(ctx).WithError(err).Errorf(errGetSession, "phone number", in.GetId())
		return nil, ErrInternal
	}

	// Check if phone number needs to be verified
	if phoneNumber.Verified {
		log.Extract(ctx).WithError(err).Errorf(errAlreadyValid, "phone", phoneNumber.ID.String(), err)
		return nil, errRetryAlreadyValid
	}

	// Check if enough time has passed to retry
	if dur := time.Until(phoneNumber.UpdatedAt.Add(time.Second * variables.VerifyPhoneNumberRetryDebounce)); dur > 0 {
		return nil, grpcerrors.ErrFailedPrecondition(fmt.Sprintf(errWait, dur/time.Second))
	}

	// Update verificationType if different
	if phoneNumber.VerificationType.String() != in.VerificationType.String() {
		phoneNumber.UpdatedAt = time.Now()

		if err = e.db.PhoneNumberDB.Update(ctx, phoneNumber); err != nil {
			log.Extract(ctx).WithError(err).Errorf(errDatabaseUpdate, "phone", phoneNumber.ID.String(), err)
			return nil, ErrInternal
		}
	}

	// Get a new verification session
	token, err := e.phoneVerifier.NewPhoneVerification(phoneNumber.PhoneNumber, in.VerificationType.String())
	if err != nil {
		log.Extract(ctx).WithError(err).Errorf(errCreateToken, "phone", phoneNumber.ID.String(), err)
		return nil, ErrInternal
	}

	// Save verification session token to DB
	phoneNumber.VerificationToken = token
	phoneNumber.UpdatedAt = time.Now()

	if err = e.db.PhoneNumberDB.Update(ctx, phoneNumber); err != nil {
		log.Extract(ctx).WithError(err).Errorf(errDatabaseUpdate, "phone", phoneNumber.ID.String(), err)
		return nil, ErrInternal
	}

	return &proto.VerifyResponse{Id: in.GetId()}, nil
}
