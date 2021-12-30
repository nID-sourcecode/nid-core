Dit hoofdstuk beschrijft de verschillende geïmplementeerde componenten voor het TWi platform. Er wordt hier onderscheid gemaakt in backend service, frontend services en envoy filters. Daarnaast wordt nog een tweetal services beschreven die gebruikt worden om de functionaliteit van het platform te laten zien in demo’s.

## Backend

- Auditlog
  - **Type**: Go gRPC service
  - **Function**: Ontvangt requests die worden afgevangen in het auditlog envoy filter. De ontvangen - informatie bevat de claims uit de JWT token, requested URL, request body en de HTTP method.
  - **Lives**: TWI namespace
- OpenID authenticatie en autorisatie service
  - **Type**: Go gRPC service
  - **Function**: OAuth / OpenID implementatie voor autoriseren van data verzoeken en het inloggen van burgers.
  - **Lives**: TWI namespace
- Auth-GQL
  - **Type**: Go GraphQL service
  - **Function**: Stelt CRUD functionaliteit bloot die gebruikt wordt door het dashboard. Hiermee wordt o.a. informatie over beschikbare scopes en beschikbare bronnen binnen het platform opgehaald.
  - **Lives**: TWI namespace
- Autopseudo
  - **Type**: Go gRPC service
  - **Function**: De auto pseudo service ontvangt requests van de auto pseudo envoy filter. Hierin worden zowel de JWT token als de request body meegestuurd. De auto pseudo service heeft toegang tot de RSA private key van de namespace waarin deze draait. De private key wordt gebruikt om de pseudoniem in de JWT token te ontsleutelen. Vervolgens wordt de $$nID:subject$$ variabele vervangen met de ontsleutelde pseudoniem en terug gestuurd naar het auto pseudo filter, die de request body vervangt.
  - **Lives**: In iedere namespace
- Dashboard
  - **Type**: Go gRPC service
  - **Function**: Basis functionaliteit voor het dashboard. O.a. inloggen van dashboard gebruikers en ophalen van informatie over welke services er draaien binnen het platform.
  - **Lives**: TWI namespace
- Documentation
  - **Type**: Go gRPC service
  - **Function**: Ophalen van documentatie die wordt weergeven in het dashboard.
  - **Lives**: TWI namespace
- JWKS
  - **Type**: Go gRPC service
  - **Function**: exposen van publieke key van RSA key pairs waarmee JWT tokens worden getekend.
  - **Lives**: TWI namespace
- Pseudonymization
  - **Type**: Go gRPC service
  - **Function**: Genereert en vertaalt pseudoniemen. Een service kan een nieuw pseudoniem aanvragen voor nieuwe gebruikers of bestaande pseudoniemen laten vertalen om te gebruiken in de communicatie met andere services. De OpenID authenticatie en autorisatie service gebruikt deze versleutelde pseudoniemen ook voor het genereren van JWT tokens.
  - **Lives**: TWI namespace
- ScopeVerification
  - **Type**: Go gRPC service
  - **Function**: Het verifiëren van requests. Hierin wordt voor graphql requests bekeken of de Graphql query matcht met de JWT token en voor normale requests of het path, body, method en query parameters matchen met het JWT token. Verder wordt gekeken of de signature van de JWTs klopt.
  - **Lives**: TWI namespace
- Wallet
  - **Type**: Go gRPC service
  - **Function**: Het registreren van nieuwe gebruikers en het inloggen in de Wallet App.
  - **Lives**: TWI namespace
- Wallet-GQL
  - **Type**: Go GraphQL service
  - **Function**: Stelt CRUD functionaliteit bloot voor informatie over de burger welke getoond wordt in de app.
  - **Lives**: TWI namespace

## Frontend

- Authrequest
  - Frontend implementatie van OAuth / OpenID. Hier komt de burger uit na het aanklikken van de nID Button. De frontend toont welke partij wat voor data wil ophalen. De burger kan vervolgens middels de Wallet App de getoonde QR code scannen en besluiten de afnemende partij toegang te verschaffen.
- Admin
  - Het dashboard gebruikt door afnemende partijen om inzage te krijgen in beschikbare scopes, draaiende services en andere vormen van documentatie.
- Wallet App
  - De App waarin de burger wordt geauthenticeerd. Deze app dient vervolgens als het middel om regie te voeren over zijn of haar data en inzage te verschaffen aan de burger welke partijen wat voor data ophalen.

## Envoy filters

- Auditlog
  - Auditlog filter leeft in de sidecar voor elke databron of informatie service en luistert naar inkomend verkeer. Elke HTTP request wordt onderschept en doorgestuurd naar de auditlog service welke de betreffende request logt. Zowel de request als responses worden doorgestuurd, echter wordt voor de response alleen de status headers gelogd. Alle logs voor eenzelfde request kunnen worden gerelateerd middels de x-request-id header.
- Authswap
  - Authswap filter is een filter dat luistert naar uitgaand verkeer op elke informatie service. Wanneer een uitgaande request wordt gedetecteerd wordt de initiële JWT token omgewisseld voor een JWT token die gebruikt kan worden voor de databron in kwestie. Het filter vervangt de Authorization header.
- Autopseudo
  - Auto Pseudo filter zit voor elke informatie service en databron en luistert naar inkomend verkeer. Wanneer een request wordt onderschept, wordt zowel de JWT token als de request body/request parameters doorgestuurd naar de auto pseudo service. Hier wordt de $$nID:subject$$ variabele vervangen en de nieuwe request body of parameters terug gestuurd. Het filter vervangt deze en stuurt de request door naar de betreffende service.
- Scopeverification
  - Scope verification filter leeft voor elke informatie service en databron en luistert naar inkomend verkeer. De scope verification stuurt alle request details naar de scope verification service die kijkt of de request mag worden gedaan. Indien de request moet worden afgewezen stuurt het filter de request met reden terug. Indien de request wordt goedgekeurd wordt deze doorgestuurd naar de betreffende service.

## Demo services

- Information
  - **Type**: Go REST Service
  - **Function**: Demo purposes om gedrag van een informatie service na te bootsen
- Databron
  - **Type**: Go GraphQL Service
  - **Function**: Demo purposes om gedrag van een databron service na te boots`
