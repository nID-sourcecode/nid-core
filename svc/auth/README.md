# Authorization Service

## Client Credentials flow

Met deze flow kunnen aangesloten clients een token opvragen voor de beschikbare scopes/audiences die geregistreerd staan in het systeem. Deze flow is vooral bedoeld voor Machine-to-Machine communicatie.

### Aanvragen van een access token

Om een access token aan te vragen, moet een POST request worden gedaan naar auth.{{host}}/token. In de body van het request moeten de volgende parameters worden meegegeven:

| Parameter  | Omschrijving                                   |
| ---------- | ---------------------------------------------- |
| grant_type | Moet de waarde 'client_credentials' hebben     |
| scopes     | Een space seperated lijst van benodigde scopes |

### Genereren van een access token

#### Audience

> **Note:** Audience claim is string die de ontvangers identificeert waarvoor de JWT bedoeld is.

Er kan maar 1 audience in een token staan. Deze audience wordt gezet op basis van de gevraagde scopes.

##### Audience Config

Door middel van de config optie `ALLOW_MULTIPLE_AUDIENCES` kan er aangegeven worden dat er meerdere audiences in een token mogen komen. Bij default is deze _false_.

#### Scopes

> **Note:** Scope is een parameter die wordt gebruikt om de reikwijdte van de toegang op te geven die wordt aangevraagd.

De scopes die gebruikt kunnen worden staan beschreven in de API specificatie.

Elke scope is aan 1 of meerdere audiences gekoppeld. Het tabel hieronder geeft een overzicht van welke response er verwacht kan worden bij het aanvragen van een of scopes op basis van configuratie.

| ALLOW_MULTIPLE_AUDIENCES | Scope audience count | Response        |
| ------------------------ | -------------------- | --------------- |
| (default) false          | 1                    | 200 OK          |
| (default) false          | > 1                  | 400 Bad Request |
| true                     | 1                    | 200 OK          |
| true                     | > 1                  | 200 OK          |

Wanneer er maar 1 audience in een token mag staan, en er worden meerdere scopes meegegeven, dan moet de client bepalen welke audience er in de token komt te staan. Dit kan door middel van de `audience` parameter in de body van het request.

#### Subject

> **Note:** Subject claim is een unieke identifier voor de gebruiker van de token.

De subject claim wordt gezet op basis van de client id die wordt meegegeven in de Authorization header.

## Dependencies
- wallet-rpc

