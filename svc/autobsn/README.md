# Automatic pseudonym to BSN translator

This service is responsible for two things:
- Decrypting pseudonyms and translating them to BSNs before arrival(`/replacePlaceholderWithBSN`)

## Decrypt and apply

`/replacePlaceholderWithBSN` is called using POST with e.g. a body:
```json
{
    "body": "the request body",
    "query": "fruit=apple&vehicle=airplane",
    "method": "POST",
    "authorization_header": "Bearer jwt.token.here"
}
```

The relevant pseudonym is extracted from the jwt claims field `sub`. E.g. using the following claims:
```json
{
  "sub" : "encrypted_system_pseudo"
} 
```

yields `encrypted_system_pseudo`. This is then decrypted, is exchanged for the BSN by calling the wallet, and any occurrence of
`$$nid:bsn$$` in the request is replaced with the BSN.

The new body and query string are then returned and can be forwarded to the actual service.
