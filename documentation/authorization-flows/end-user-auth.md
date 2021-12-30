# End-user authorzation
This document describes the process of end-user authorization; when and end-user grants consent for a _client_ to access some of the end-user's data supplied by an _audience_.

The end-user authorization flow of the authorization component adheres to the [openID spec](https://openid.net/specs/openid-connect-core-1_0.html). In short:

```plantuml
actor Client as client
actor Audience as audience
actor "End-user" as user
participant "authrequest frontend" as authrequest
participant "auth service" as auth
participant "psuedonymization service" as pseudo

client->user : send /authorize url
user->auth : /authorize
hnote over auth: authenticate user
auth->user: redirect to client's redirect-uri
user->client: redirected to redirect uri
client->user: show landing page
client->auth: retrieve token
auth->pseudo: translate psuedonym
pseudo->auth: translated psuedonym
hnote over auth: creating token
auth->client: send token response
client->audience: get data using token
audience->client: data
```

The 'authenticate user' step looks as follows:

```plantuml
actor "End-user" as user
participant "authrequest frontend" as authrequest
participant "auth service" as auth

auth->user: redirect to authrequest
user->authrequest : redirected
authrequest->auth : get info
auth->authrequest: send info
authrequest-->auth : poll session status
user->authrequest : scan QR code using phone
user->auth : claim session based on QR
auth-->authrequest: notify session claimed
user->auth : grant access on phone
auth-->authrequest: notify access granted
user->authrequest: confirm access granted
authrequest->user: redirect
user->auth: redirected to /finalize
```

## Flow of sensitive information
```plantuml
actor Client as client
actor Audience as audience
actor "End-user" as user
participant "authrequest frontend" as authrequest
participant "auth service" as auth
participant "psuedonymization service" as pseudo

client->user: audience + scopes<sup>1</sup>
user->auth: audience + scopes<sup>1</sup>
auth->authrequest: audience + <sup>2</sup>
user->auth: wallet pseudonym<sup>3</sup>
auth->user: auth token<sup>4</sup>
user->client: auth token<sup>4</sup>
client->auth: client password<sup>5</sup>
auth->pseudo: wallet pseudonym<sup>6</sup>
pseudo->auth: encrypted audience psuedo<sup>7</sup>
auth->client: JWT token containing encrypted pseudos and scopes<sup>8</sup>
client->audience: JWT<sup>9</sup>
audience->audience: decrypted audience pseudo<sup>10</sup>
audience->client: Data<sup>11</sup>
```

<!-- FIXME add references to the swagger here -->
#### Notes
1. Encoded in the /authorize url parameters
2. Part of the /getsession response
3. Present in the JWT received from the wallet in the auth header with the /claim call
4. Encoded in the client redirect uri parameters.
5. Present in basic auth for the /token call from the client.
6. Request for /Convert call.
7. Response for /Convert call.
8. Response for /token call.
9. Auth header for data call.
10. Decrypted by `autopseudo` service in audience namespace (wasm filter).
11. Response for data call.