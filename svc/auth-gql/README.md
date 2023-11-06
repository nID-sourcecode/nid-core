# Auth GQL service
The auth GQL service exposes a GraphQL API which allows to perform the following actions:
 - Create or refresh a JWT
 - Query for Client(s)

## Development
This service generates the graphql schema from the model design inside the auth service. After modifying the design inside `auth/design`, the schema can be generated with:
```bash
cd ./svc/auth
gen -graphqlPath=../auth-gql/graphql -authPath=../auth-gql/auth
```