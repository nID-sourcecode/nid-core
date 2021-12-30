In dit document wordt er uitgeweid over het aansluiten van zowel afnemende partijen als databronnen. Ook burgers zullen een vorm van aansluiten nodig gaan hebben. Hier wordt momenteel gedacht aan DigiD, maar moet nog verder uitgewerkt worden in PI2 & 3.

## Aansluiten Afnemende partij

Momenteel is het mogelijk om afnemende partijen handmatig aan te sluiten op het platform. Voor een afnemende partijen worden de volgende attributen aangemaakt voor zowel een staging als productie omgeving:

- Username & password om toegang te geven tot het dashboard
- Een clientID

Door het pilot team moet worden aangeleverd;

- Een redirect URI
- Een JWKS endpoint

Verder wordt aan het pilot team aangeleverd:

- Platform URL
- Documentatie middels het dashboard
- nID Button
- nID Wallet app

In PI3 zal er gewerkt moeten worden aan een (semi) geautomatiseerd vetting process waarmee nieuwe afnemende partijen ge-onboard kunnen worden.

## Aansluiten Databron

Wanneer een nieuwe databron wordt aangesloten wordt er voor de betreffende databron een namespace aangemaakt. Vervolgens kan de databron op 2 manieren aansluiten. Allereerst kan er gezorgd worden dat de databron zijn services direct deployed binnen het platform. De tweede optie is dat de databron middels [mesh expansion](https://istio.io/v1.1/docs/setup/kubernetes/additional-setup/mesh-expansion/) wordt aangesloten en services in zijn eigen omgeving kan blijven beheren. Aangezien het platform werkt met pseudoniemen en bronnen veelal nog zullen werken met BSN’s zal er tijdens het onboarding proces eenmalig een translatie slag gemaakt moeten worden door een bron. Dit betekend dat voor elk BSN, een databron een pseudoniem moet aanvragen binnen het platform. Vervolgens kan de databron de data relateren aan de pseudoniemen in plaats van aan BSN’s. Mocht de databron de data toch nog willen blijven relateren aan BSN’s voor andere doeleinden, kan deze informatie bijvoorbeeld in een offline omgeving bewaard blijven.

Voor de eerst opgeleverde versie van PI1 zal er voor BKWI deze translatie slag gemaakt worden aan de kant van het platform, om voor nu BKWI te ontzien van extra werk. De bedoeling is dat wanneer de GraphQL services klaar zijn, BKWI deze translatieslag zelf gaat beheren.

## Informatieservices

Een afnemende partij vraagt aan een burger om toestemming tot een specifieke informatie service. Deze informatie service heeft 1 of meerdere bronnen nodig om antwoord te kunnen geven op een vraag. Hiervoor moet de informatie service de mogelijkheid hebben om databronnen te bevragen. Om dit te realiseren, is er voor informatieservices een uitgaan envoy filter toegevoegd die het initiële JWT token ingewisseld voor een JWT token specifiek voor de opgevraagde databron. Dit werkt middels een intern endpoint binnen het platform die alleen toegankelijk is door het envoy filter van informatieservices. In de eerste versie is een informatie service in staat om zijn initiële JWT token in te wisselen voor elke databron die is aangesloten op het platform. In PI2 zal er een rollen systeem toegevoegd worden waardoor een informatie service, alsmede een afnemende partij, restricties krijgt tot welke bronnen zij informatie op mogen vragen.
