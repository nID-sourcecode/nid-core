# Wallet GQL service
The wallet GQL service exposes a GraphQL API which allows to perform the following actions:
 - Query for client(s)
 - Update / create client(s)
 - Query for consent(s)
 - Revoke consent
 - Query for emailAddress(es)
 - Update / Create emailAddress(es)
 - Query for phoneNumber(s)
 - Update / Create phoneNumber(s)
 - Query for user(s)
 - Create / Update user(s)
 - Create or refresh a JWT
 

## Development
This service generates the graphql schema from the model design inside the wallet-gql service. After modifying the design inside `wallet-gql/design`, the schema can be generated with:
```bash
cd ./svc/wallet-gql
gen
```