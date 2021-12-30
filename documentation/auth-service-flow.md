# Auth service flow

## Authorization url received

The registered client sends an authorization url to the user (by email for example). The user clicks on this authorization url whereafter a new session for this user will be created. After the session is created the user redirects to the configured user authentication frontend which in our case is the authrequest web-application.

**Endpoint:** /authorize  
**Type:** GET  
**Parameters:** client_id, response_type, scope, redirect_uri, audience, optional_scopes  
**Status code:** 302  
**Response body:**  

## Receiving the access token

When the finalise action delivers the authorization code to the client. The client can exchange this authorization code for an access token. This way the client gets access to fetch user data.

**Endpoint:** /token  
**Type:** GET  
**Parameters:**  grant_type, authorization_code  
**Status code:** 200  
**Response body:**  access_token, token_type
