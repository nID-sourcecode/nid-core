// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
)

type BooleanFilterInput struct {
	Eq    *bool               `json:"eq"`
	Ne    *bool               `json:"ne"`
	IsSet *bool               `json:"isSet"`
	Not   *BooleanFilterInput `json:"not"`
}

type ClientFilterInput struct {
	ID          *UUIDFilterInput     `json:"id"`
	Color       *StringFilterInput   `json:"color"`
	ExtClientID *StringFilterInput   `json:"extClientId"`
	Icon        *StringFilterInput   `json:"icon"`
	Logo        *StringFilterInput   `json:"logo"`
	Name        *StringFilterInput   `json:"name"`
	CreatedAt   *TimeFilterInput     `json:"createdAt"`
	DeletedAt   *TimeFilterInput     `json:"deletedAt"`
	UpdatedAt   *TimeFilterInput     `json:"updatedAt"`
	Unscoped    *bool                `json:"unscoped"`
	Not         *ClientFilterInput   `json:"not"`
	Or          []*ClientFilterInput `json:"or"`
	And         []*ClientFilterInput `json:"and"`
}

type ConsentFilterInput struct {
	ID          *UUIDFilterInput      `json:"id"`
	AccessToken *StringFilterInput    `json:"accessToken"`
	ClientID    *UUIDFilterInput      `json:"clientId"`
	Description *StringFilterInput    `json:"description"`
	Granted     *TimeFilterInput      `json:"granted"`
	Name        *StringFilterInput    `json:"name"`
	Revoked     *TimeFilterInput      `json:"revoked"`
	Token       *JSONFilterInput      `json:"token"`
	UserID      *UUIDFilterInput      `json:"userId"`
	CreatedAt   *TimeFilterInput      `json:"createdAt"`
	DeletedAt   *TimeFilterInput      `json:"deletedAt"`
	UpdatedAt   *TimeFilterInput      `json:"updatedAt"`
	Unscoped    *bool                 `json:"unscoped"`
	Not         *ConsentFilterInput   `json:"not"`
	Or          []*ConsentFilterInput `json:"or"`
	And         []*ConsentFilterInput `json:"and"`
}

type CreateClient struct {
	Color       string `json:"color"`
	ExtClientID string `json:"extClientId"`
	Icon        string `json:"icon"`
	Logo        string `json:"logo"`
	Name        string `json:"name"`
}

type CreateEmailAddress struct {
	EmailAddress string     `json:"emailAddress"`
	UserID       *uuid.UUID `json:"userId"`
}

type CreatePhoneNumber struct {
	PhoneNumber      string                             `json:"phoneNumber"`
	UserID           *uuid.UUID                         `json:"userId"`
	VerificationType models.PhoneNumberVerificationType `json:"verificationType"`
}

type CreateRevokeConsent struct {
	ID *uuid.UUID `json:"id"`
}

type CreateUser struct {
	Bsn       string          `json:"bsn"`
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	Pseudonym string          `json:"pseudonym"`
	Scopes    *postgres.Jsonb `json:"scopes"`
}

type EmailAddressFilterInput struct {
	ID           *UUIDFilterInput           `json:"id"`
	EmailAddress *StringFilterInput         `json:"emailAddress"`
	UserID       *UUIDFilterInput           `json:"userId"`
	Verified     *BooleanFilterInput        `json:"verified"`
	CreatedAt    *TimeFilterInput           `json:"createdAt"`
	DeletedAt    *TimeFilterInput           `json:"deletedAt"`
	UpdatedAt    *TimeFilterInput           `json:"updatedAt"`
	Unscoped     *bool                      `json:"unscoped"`
	Not          *EmailAddressFilterInput   `json:"not"`
	Or           []*EmailAddressFilterInput `json:"or"`
	And          []*EmailAddressFilterInput `json:"and"`
}

type JSONFilterInput struct {
	Contains *postgres.Jsonb  `json:"contains"`
	Eq       *postgres.Jsonb  `json:"eq"`
	HasPath  []string         `json:"hasPath"`
	Ne       *postgres.Jsonb  `json:"ne"`
	IsSet    *bool            `json:"isSet"`
	Not      *JSONFilterInput `json:"not"`
}

type PhoneNumberFilterInput struct {
	ID               *UUIDFilterInput                        `json:"id"`
	PhoneNumber      *StringFilterInput                      `json:"phoneNumber"`
	UserID           *UUIDFilterInput                        `json:"userId"`
	VerificationType *PhoneNumberVerificationTypeFilterInput `json:"verificationType"`
	Verified         *BooleanFilterInput                     `json:"verified"`
	CreatedAt        *TimeFilterInput                        `json:"createdAt"`
	DeletedAt        *TimeFilterInput                        `json:"deletedAt"`
	UpdatedAt        *TimeFilterInput                        `json:"updatedAt"`
	Unscoped         *bool                                   `json:"unscoped"`
	Not              *PhoneNumberFilterInput                 `json:"not"`
	Or               []*PhoneNumberFilterInput               `json:"or"`
	And              []*PhoneNumberFilterInput               `json:"and"`
}

type PhoneNumberVerificationTypeFilterInput struct {
	Eq    *models.PhoneNumberVerificationType     `json:"eq"`
	In    []models.PhoneNumberVerificationType    `json:"in"`
	Ne    *models.PhoneNumberVerificationType     `json:"ne"`
	IsSet *bool                                   `json:"isSet"`
	Not   *PhoneNumberVerificationTypeFilterInput `json:"not"`
}

type RevokeConsent struct {
	ID      uuid.UUID `json:"id"`
	Revoked time.Time `json:"revoked"`
}

type StringFilterInput struct {
	BeginsWith *string            `json:"beginsWith"`
	Contains   *string            `json:"contains"`
	EndsWith   *string            `json:"endsWith"`
	Eq         *string            `json:"eq"`
	Ge         *string            `json:"ge"`
	Gt         *string            `json:"gt"`
	Le         *string            `json:"le"`
	Lt         *string            `json:"lt"`
	Ne         *string            `json:"ne"`
	IsSet      *bool              `json:"isSet"`
	Not        *StringFilterInput `json:"not"`
}

type TimeFilterInput struct {
	Eq    *time.Time       `json:"eq"`
	Ge    *time.Time       `json:"ge"`
	Gt    *time.Time       `json:"gt"`
	Le    *time.Time       `json:"le"`
	Lt    *time.Time       `json:"lt"`
	Ne    *time.Time       `json:"ne"`
	IsSet *bool            `json:"isSet"`
	Not   *TimeFilterInput `json:"not"`
}

type UUIDFilterInput struct {
	Eq    *uuid.UUID       `json:"eq"`
	Ne    *uuid.UUID       `json:"ne"`
	IsSet *bool            `json:"isSet"`
	Not   *UUIDFilterInput `json:"not"`
}

type UpdateClient struct {
	Color       *string `json:"color"`
	ExtClientID *string `json:"extClientId"`
	Icon        *string `json:"icon"`
	Logo        *string `json:"logo"`
	Name        *string `json:"name"`
}

type UpdateEmailAddress struct {
	EmailAddress *string    `json:"emailAddress"`
	UserID       *uuid.UUID `json:"userId"`
}

type UpdatePhoneNumber struct {
	PhoneNumber      *string                             `json:"phoneNumber"`
	UserID           *uuid.UUID                          `json:"userId"`
	VerificationType *models.PhoneNumberVerificationType `json:"verificationType"`
}

type UpdateUser struct {
	Bsn       *string         `json:"bsn"`
	Email     *string         `json:"email"`
	Password  *string         `json:"password"`
	Pseudonym *string         `json:"pseudonym"`
	Scopes    *postgres.Jsonb `json:"scopes"`
}

type UserFilterInput struct {
	ID        *UUIDFilterInput   `json:"id"`
	Bsn       *StringFilterInput `json:"bsn"`
	Email     *StringFilterInput `json:"email"`
	Pseudonym *StringFilterInput `json:"pseudonym"`
	Scopes    *JSONFilterInput   `json:"scopes"`
	CreatedAt *TimeFilterInput   `json:"createdAt"`
	DeletedAt *TimeFilterInput   `json:"deletedAt"`
	UpdatedAt *TimeFilterInput   `json:"updatedAt"`
	Unscoped  *bool              `json:"unscoped"`
	Not       *UserFilterInput   `json:"not"`
	Or        []*UserFilterInput `json:"or"`
	And       []*UserFilterInput `json:"and"`
}

type ClientFieldName string

const (
	ClientFieldNameID          ClientFieldName = "ID"
	ClientFieldNameColor       ClientFieldName = "COLOR"
	ClientFieldNameExtClientID ClientFieldName = "EXT_CLIENT_ID"
	ClientFieldNameIcon        ClientFieldName = "ICON"
	ClientFieldNameLogo        ClientFieldName = "LOGO"
	ClientFieldNameName        ClientFieldName = "NAME"
	ClientFieldNameCreatedAt   ClientFieldName = "CREATED_AT"
	ClientFieldNameDeletedAt   ClientFieldName = "DELETED_AT"
	ClientFieldNameUpdatedAt   ClientFieldName = "UPDATED_AT"
)

var AllClientFieldName = []ClientFieldName{
	ClientFieldNameID,
	ClientFieldNameColor,
	ClientFieldNameExtClientID,
	ClientFieldNameIcon,
	ClientFieldNameLogo,
	ClientFieldNameName,
	ClientFieldNameCreatedAt,
	ClientFieldNameDeletedAt,
	ClientFieldNameUpdatedAt,
}

func (e ClientFieldName) IsValid() bool {
	switch e {
	case ClientFieldNameID, ClientFieldNameColor, ClientFieldNameExtClientID, ClientFieldNameIcon, ClientFieldNameLogo, ClientFieldNameName, ClientFieldNameCreatedAt, ClientFieldNameDeletedAt, ClientFieldNameUpdatedAt:
		return true
	}
	return false
}

func (e ClientFieldName) String() string {
	return string(e)
}

func (e *ClientFieldName) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ClientFieldName(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ClientFieldName", str)
	}
	return nil
}

func (e ClientFieldName) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ConsentFieldName string

const (
	ConsentFieldNameID          ConsentFieldName = "ID"
	ConsentFieldNameAccessToken ConsentFieldName = "ACCESS_TOKEN"
	ConsentFieldNameClientID    ConsentFieldName = "CLIENT_ID"
	ConsentFieldNameDescription ConsentFieldName = "DESCRIPTION"
	ConsentFieldNameGranted     ConsentFieldName = "GRANTED"
	ConsentFieldNameName        ConsentFieldName = "NAME"
	ConsentFieldNameRevoked     ConsentFieldName = "REVOKED"
	ConsentFieldNameToken       ConsentFieldName = "TOKEN"
	ConsentFieldNameUserID      ConsentFieldName = "USER_ID"
	ConsentFieldNameCreatedAt   ConsentFieldName = "CREATED_AT"
	ConsentFieldNameDeletedAt   ConsentFieldName = "DELETED_AT"
	ConsentFieldNameUpdatedAt   ConsentFieldName = "UPDATED_AT"
)

var AllConsentFieldName = []ConsentFieldName{
	ConsentFieldNameID,
	ConsentFieldNameAccessToken,
	ConsentFieldNameClientID,
	ConsentFieldNameDescription,
	ConsentFieldNameGranted,
	ConsentFieldNameName,
	ConsentFieldNameRevoked,
	ConsentFieldNameToken,
	ConsentFieldNameUserID,
	ConsentFieldNameCreatedAt,
	ConsentFieldNameDeletedAt,
	ConsentFieldNameUpdatedAt,
}

func (e ConsentFieldName) IsValid() bool {
	switch e {
	case ConsentFieldNameID, ConsentFieldNameAccessToken, ConsentFieldNameClientID, ConsentFieldNameDescription, ConsentFieldNameGranted, ConsentFieldNameName, ConsentFieldNameRevoked, ConsentFieldNameToken, ConsentFieldNameUserID, ConsentFieldNameCreatedAt, ConsentFieldNameDeletedAt, ConsentFieldNameUpdatedAt:
		return true
	}
	return false
}

func (e ConsentFieldName) String() string {
	return string(e)
}

func (e *ConsentFieldName) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ConsentFieldName(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ConsentFieldName", str)
	}
	return nil
}

func (e ConsentFieldName) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type EmailAddressFieldName string

const (
	EmailAddressFieldNameID                EmailAddressFieldName = "ID"
	EmailAddressFieldNameEmailAddress      EmailAddressFieldName = "EMAIL_ADDRESS"
	EmailAddressFieldNameUserID            EmailAddressFieldName = "USER_ID"
	EmailAddressFieldNameVerificationToken EmailAddressFieldName = "VERIFICATION_TOKEN"
	EmailAddressFieldNameVerified          EmailAddressFieldName = "VERIFIED"
	EmailAddressFieldNameCreatedAt         EmailAddressFieldName = "CREATED_AT"
	EmailAddressFieldNameDeletedAt         EmailAddressFieldName = "DELETED_AT"
	EmailAddressFieldNameUpdatedAt         EmailAddressFieldName = "UPDATED_AT"
)

var AllEmailAddressFieldName = []EmailAddressFieldName{
	EmailAddressFieldNameID,
	EmailAddressFieldNameEmailAddress,
	EmailAddressFieldNameUserID,
	EmailAddressFieldNameVerificationToken,
	EmailAddressFieldNameVerified,
	EmailAddressFieldNameCreatedAt,
	EmailAddressFieldNameDeletedAt,
	EmailAddressFieldNameUpdatedAt,
}

func (e EmailAddressFieldName) IsValid() bool {
	switch e {
	case EmailAddressFieldNameID, EmailAddressFieldNameEmailAddress, EmailAddressFieldNameUserID, EmailAddressFieldNameVerificationToken, EmailAddressFieldNameVerified, EmailAddressFieldNameCreatedAt, EmailAddressFieldNameDeletedAt, EmailAddressFieldNameUpdatedAt:
		return true
	}
	return false
}

func (e EmailAddressFieldName) String() string {
	return string(e)
}

func (e *EmailAddressFieldName) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EmailAddressFieldName(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EmailAddressFieldName", str)
	}
	return nil
}

func (e EmailAddressFieldName) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OrderDirection string

const (
	OrderDirectionAsc  OrderDirection = "ASC"
	OrderDirectionDesc OrderDirection = "DESC"
)

var AllOrderDirection = []OrderDirection{
	OrderDirectionAsc,
	OrderDirectionDesc,
}

func (e OrderDirection) IsValid() bool {
	switch e {
	case OrderDirectionAsc, OrderDirectionDesc:
		return true
	}
	return false
}

func (e OrderDirection) String() string {
	return string(e)
}

func (e *OrderDirection) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderDirection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderDirection", str)
	}
	return nil
}

func (e OrderDirection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PhoneNumberFieldName string

const (
	PhoneNumberFieldNameID                PhoneNumberFieldName = "ID"
	PhoneNumberFieldNamePhoneNumber       PhoneNumberFieldName = "PHONE_NUMBER"
	PhoneNumberFieldNameUserID            PhoneNumberFieldName = "USER_ID"
	PhoneNumberFieldNameVerificationToken PhoneNumberFieldName = "VERIFICATION_TOKEN"
	PhoneNumberFieldNameVerificationType  PhoneNumberFieldName = "VERIFICATION_TYPE"
	PhoneNumberFieldNameVerified          PhoneNumberFieldName = "VERIFIED"
	PhoneNumberFieldNameCreatedAt         PhoneNumberFieldName = "CREATED_AT"
	PhoneNumberFieldNameDeletedAt         PhoneNumberFieldName = "DELETED_AT"
	PhoneNumberFieldNameUpdatedAt         PhoneNumberFieldName = "UPDATED_AT"
)

var AllPhoneNumberFieldName = []PhoneNumberFieldName{
	PhoneNumberFieldNameID,
	PhoneNumberFieldNamePhoneNumber,
	PhoneNumberFieldNameUserID,
	PhoneNumberFieldNameVerificationToken,
	PhoneNumberFieldNameVerificationType,
	PhoneNumberFieldNameVerified,
	PhoneNumberFieldNameCreatedAt,
	PhoneNumberFieldNameDeletedAt,
	PhoneNumberFieldNameUpdatedAt,
}

func (e PhoneNumberFieldName) IsValid() bool {
	switch e {
	case PhoneNumberFieldNameID, PhoneNumberFieldNamePhoneNumber, PhoneNumberFieldNameUserID, PhoneNumberFieldNameVerificationToken, PhoneNumberFieldNameVerificationType, PhoneNumberFieldNameVerified, PhoneNumberFieldNameCreatedAt, PhoneNumberFieldNameDeletedAt, PhoneNumberFieldNameUpdatedAt:
		return true
	}
	return false
}

func (e PhoneNumberFieldName) String() string {
	return string(e)
}

func (e *PhoneNumberFieldName) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PhoneNumberFieldName(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PhoneNumberFieldName", str)
	}
	return nil
}

func (e PhoneNumberFieldName) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RevokeConsentFieldName string

const (
	RevokeConsentFieldNameID      RevokeConsentFieldName = "ID"
	RevokeConsentFieldNameRevoked RevokeConsentFieldName = "REVOKED"
)

var AllRevokeConsentFieldName = []RevokeConsentFieldName{
	RevokeConsentFieldNameID,
	RevokeConsentFieldNameRevoked,
}

func (e RevokeConsentFieldName) IsValid() bool {
	switch e {
	case RevokeConsentFieldNameID, RevokeConsentFieldNameRevoked:
		return true
	}
	return false
}

func (e RevokeConsentFieldName) String() string {
	return string(e)
}

func (e *RevokeConsentFieldName) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RevokeConsentFieldName(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RevokeConsentFieldName", str)
	}
	return nil
}

func (e RevokeConsentFieldName) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UserFieldName string

const (
	UserFieldNameID        UserFieldName = "ID"
	UserFieldNameBsn       UserFieldName = "BSN"
	UserFieldNameEmail     UserFieldName = "EMAIL"
	UserFieldNamePassword  UserFieldName = "PASSWORD"
	UserFieldNamePseudonym UserFieldName = "PSEUDONYM"
	UserFieldNameScopes    UserFieldName = "SCOPES"
	UserFieldNameCreatedAt UserFieldName = "CREATED_AT"
	UserFieldNameDeletedAt UserFieldName = "DELETED_AT"
	UserFieldNameUpdatedAt UserFieldName = "UPDATED_AT"
)

var AllUserFieldName = []UserFieldName{
	UserFieldNameID,
	UserFieldNameBsn,
	UserFieldNameEmail,
	UserFieldNamePassword,
	UserFieldNamePseudonym,
	UserFieldNameScopes,
	UserFieldNameCreatedAt,
	UserFieldNameDeletedAt,
	UserFieldNameUpdatedAt,
}

func (e UserFieldName) IsValid() bool {
	switch e {
	case UserFieldNameID, UserFieldNameBsn, UserFieldNameEmail, UserFieldNamePassword, UserFieldNamePseudonym, UserFieldNameScopes, UserFieldNameCreatedAt, UserFieldNameDeletedAt, UserFieldNameUpdatedAt:
		return true
	}
	return false
}

func (e UserFieldName) String() string {
	return string(e)
}

func (e *UserFieldName) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserFieldName(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserFieldName", str)
	}
	return nil
}

func (e UserFieldName) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
