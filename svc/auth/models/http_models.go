package models

// AuthorizeRequest holds the parameters for the authorization process.
type AuthorizeRequest struct {
	Scope          string `json:"scope,omitempty"`
	ResponseType   string `json:"response_type,omitempty"`
	ClientID       string `json:"client_id,omitempty"`
	RedirectURI    string `json:"redirect_uri,omitempty"`
	Audience       string `json:"audience,omitempty"`
	OptionalScopes string `json:"optional_scopes,omitempty"`
}

func (a *AuthorizeRequest) GetClientID() string {
	return a.ClientID
}

func (a *AuthorizeRequest) GetAudience() string {
	return a.Audience
}

func (a *AuthorizeRequest) GetRedirectURI() string {
	return a.RedirectURI
}

// AuthorizeHeadlessRequest is used in a headless authorization flow.
type AuthorizeHeadlessRequest struct {
	ResponseType   string `json:"response_type,omitempty"`
	ClientID       string `json:"client_id,omitempty"`
	RedirectURI    string `json:"redirect_uri,omitempty"`
	Audience       string `json:"audience,omitempty"`
	QueryModelJSON string `json:"query_model_json,omitempty"`
	QueryModelPath string `json:"query_model_path,omitempty"`
}

func (a *AuthorizeHeadlessRequest) GetClientID() string {
	return a.ClientID
}

func (a *AuthorizeHeadlessRequest) GetAudience() string {
	return a.Audience
}

func (a *AuthorizeHeadlessRequest) GetRedirectURI() string {
	return a.RedirectURI
}

// SessionRequest holds the unique identifier of a session.
type SessionRequest struct {
	SessionID string `json:"session_id,omitempty"`
}

// AcceptRequest is used to approve access to certain resources.
type AcceptRequest struct {
	SessionID      string   `json:"session_id,omitempty"`
	AccessModelIds []string `json:"access_model_ids,omitempty"`
}

// FinaliseRequest is used to end the session after ensuring everything is in order.
type FinaliseRequest struct {
	SessionID            string `json:"session_id,omitempty"`
	SessionFinaliseToken string `json:"session_finalise_token,omitempty"`
}

// TokenRequestMetadata is used to pass additional information to the token endpoint.
type TokenRequestMetadata struct {
	CertificateHeader string `json:"client_id,omitempty"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
}

// TokenRequest is used to request a new token or refresh an existing one.
type TokenRequest struct {
	GrantType         string               `json:"grant_type"`
	AuthorizationCode string               `json:"authorization_code,omitempty"`
	RefreshToken      string               `json:"refresh_token,omitempty"`
	Metadata          TokenRequestMetadata `json:"metadata,omitempty"`
}

// TokenClientFlowRequest is used to request a new token in a client flow.
type TokenClientFlowRequest struct {
	GrantType string               `json:"grant_type,omitempty"`
	Scope     string               `json:"scope,omitempty"`
	Audience  string               `json:"audience,omitempty"`
	Metadata  TokenRequestMetadata `json:"metadata,omitempty"`
}

// AccessModelRequest is used to request access to certain resources.
type AccessModelRequest struct {
	Audience       string `json:"audience,omitempty"`
	QueryModelJSON string `json:"query_model_json,omitempty"`
	ScopeName      string `json:"scope_name,omitempty"`
	Description    string `json:"description,omitempty"`
}

// SwapTokenRequest is used to exchange a token for another one.
type SwapTokenRequest struct {
	CurrentToken string `json:"currentToken,omitempty"`
	Query        string `json:"query,omitempty"`
	Audience     string `json:"audience,omitempty"`
}

// SessionResponse provides the state and detail of a session.
type SessionResponse struct {
	ID                   string         `json:"id,omitempty"`
	State                SessionState   `json:"state,omitempty"`
	Client               *Client        `json:"client,omitempty"`
	Audience             *Audience      `json:"audience,omitempty"`
	RequiredAccessModels []*AccessModel `json:"required_access_models,omitempty"`
	OptionalAccessModels []*AccessModel `json:"optional_access_models,omitempty"`
	AcceptedAccessModels []*AccessModel `json:"accepted_access_models,omitempty"`
}

// SessionAuthorization provides the finalization token for a session.
type SessionAuthorization struct {
	FinaliseToken string `json:"finaliseToken,omitempty"`
}

// StatusResponse provides the state of a session.
type StatusResponse struct {
	State SessionState `json:"state,omitempty"`
}

// FinaliseResponse provides the location where client should be redirected after finalizing.
type FinaliseResponse struct {
	RedirectLocation string `json:"redirect_location,omitempty"`
}

// TokenResponse provides the token details.
type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}
