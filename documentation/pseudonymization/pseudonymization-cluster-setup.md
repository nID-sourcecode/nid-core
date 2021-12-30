# State of pseudonymization in the TWI cluster

This document describes the state of pseudonymization in the TWI cluster, and what this means for onboarding parties and platform functionality.

## Table of contents
[toc]

## Platform functionality

Pseudonymization has several consequences for platform functionality.

- Access tokens contain encrypted pseudonyms instead of BSNs
- Audiences have filters that decrypt incoming encrypted pseudonyms automatically

### Granting consent
In the flow for getting an access token, pseudonymization only plays a role at the very end of the process, when creating the token. At that point, the wallet's pseudonym is translated to the audience's pseudonym.

```plantuml
actor Client as client
participant "auth service" as auth
participant "pseudonymization service" as pseudo
client -> auth: request token
auth-> pseudo: translate pseudonym
pseudo ->auth: send translated pseudonym
hnote over auth
create token
with pseudonym
as subject
end note
auth->client: send token
```

The same goes for the automatic role-based consent:
```plantuml
participant "information-svc" as client
participant "auth-svc" as auth
participant "pseudonymization service" as pseudo
client -> auth: request token based on role
auth-> pseudo: translate pseudonym
pseudo ->auth: send translated pseudonym
hnote over auth
create token
with pseudonym
as subject
end note
auth->client: send token
```

### Requesting data
```plantuml
actor Client as client
participant "Ingress" as gress
participant "autopseudo-filter" as autopseudo
participant Audience as audience
client->gress: request data from audience (JWT in header)
gress->autopseudo: request data (JWT in header)
hnote over autopseudo: decrypt pseudonym
hnote over autopseudo: replace $$nid:subject$$ with pseudonym
autopseudo->audience: request data
audience->gress: send data
gress->client: send data
```

## Onboarding
Using a pseudonymized system requires various forms of onboarding from involved parties; everyone needs to be on the same page regarding pseudonym usage. This section describes what is needed for pseudonymization to function.

### Audience onboarding
If data is to be retrieved from an audience on the basis of pseudonyms, the audience needs to know which pseudonym belongs to whom. To that end, the audience needs to translate its BSNs to pseudonyms, and store those in its database.

```plantuml
actor Audience as audience
participant "onboarding service" as onboarding
audience->onboarding: Request to convert BSNs to pseudonyms
onboarding->audience: Send mapping from BSNs to pseudonyms
hnote over audience: Save pseudonyms in an extra column on burger table
```

### End-user onboarding
```plantuml
actor "End-user" as user
participant "Wallet app" as app
participant DigiD as digid
participant "Wallet backend" as wallet
participant Pseudonymization as pseudo
user->app: sign up
app->digid: sign in with digiD
digid->wallet: bsn
group if burger does not yet exist
    wallet->pseudo: generate pseudonym
    pseudo->wallet: send new pseudonym
    hnote over wallet
    create new burger
    with received
    pseudonym and BSN
    end note
end
wallet->app: sign in as burger
```

### Client onboarding
A client is expected to know the subject of the information they are asking (by tracking a session in the authorization flow), and therefore has no need for onboarding within the scope of data authorization.

Client onboarding may become relevant when we want to provide _end-user authentication_.

## Platform functionality

Pseudonymization has several consequences for platform functionality.
- Access tokens contain encrypted pseudonyms instead of BSNs
- Audiences have filters that decrypt incoming encrypted pseudonyms automatically

### Granting consent
In the flow for getting an access token, pseudonymization only plays a role at the very end of the process, when creating the token. At that point, the wallet's pseudonym is translated to the audience's pseudonym.

```plantuml
actor Client as client
participant "auth service" as auth
participant "pseudonymization service" as pseudo
client -> auth: request token
auth-> pseudo: translate pseudonym
pseudo ->auth: send translated pseudonym
hnote over auth
create token
with pseudonym
as subject
end note
auth->client: send token
```

The same goes for the automatic role-based consent:
```plantuml
participant "information service filter" as client
participant "auth service" as auth
participant "pseudonymization service" as pseudo
client -> auth: request token based on role
auth-> pseudo: translate pseudonym
pseudo ->auth: send translated pseudonym
hnote over auth
create token
with pseudonym
as subject
end note
auth->client: send token
```

### Requesting data
```plantuml
actor Client as client
participant "Ingress" as gress
participant Autopseudo as autopseudo
participant Audience as audience
client->gress: request data from audience (JWT in header)
gress->autopseudo: request data (JWT in header)
hnote over autopseudo: decrypt pseudonym
hnote over autopseudo: replace $$nid:subject$$ with pseudonym
autopseudo->audience: request data
audience->gress: send data
gress->client: send data
```

<!--## Hybrid pseudonymization setup (proposal)
TODO do something with this. Remove or polish.

Rather than changing many aspects of the cluster to work without pseudonymization, we create filters around parties who wish to work with BSNs. This creates a sort of _hybrid_ state, where pseudonymization is _enabled_, but not _required_ for parties using the cluster.

How does this work?

```plantuml
actor Client as client
participant "Ingress" as gress
participant Autopseudo as autopseudo
participant Audience as audience
client->gress: request data from audience (JWT in header)
gress->autopseudo: request data (JWT in header)
hnote over autopseudo: decrypt pseudonym
hnote over autopseudo: translate pseudonym to BSN
hnote over autopseudo: replace $$nid:subject$$ with BSN
autopseudo->audience: request data
audience->gress: send data
gress->client: send data
```

This removes the need for **audience onboarding**.-->
