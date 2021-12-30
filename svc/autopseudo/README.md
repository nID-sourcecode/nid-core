# Automatic Pseudonym Decrypter

This service is responsible for two things:
- Storing and exposing the pseudonym encryption jwks for a namespace (`/jwks`)
- Decrypting pseudonyms before arrival (`/decryptAndApply`)

## Decrypt and apply

`/decryptAndApply` is called just like a gql endpoint. The relevant pseudonym is extracted from the jwt claims field `subject`. E.g. using the following claims:
```json
{
  "subject" : {
    "alice": "encryptedpseudo1"
  } 
}
```

yields `encryptedpseudo1`, if this service is run in the namespace alice. This is then decrypted, and any occurrence of
`$$nid:subject$$` in the query is replaced with the decrypted pseudonym.

The new gql query is then returned and can be forwarded to the actual service.