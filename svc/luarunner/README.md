# Luarunner
The Luarunner service exposes a basic HTTP api that can be used to execute Luascripts.
Currently this service implements one use case: Consent by law. This use case is tightly coupled to the `authorizeHeadless` function inside the `auth` service.

The Luarunner fetches the correct redirect URL and linked client from the organisationId inside the request

Example request to the Luarunner:
```
{
    "organisatieId" : "1",
    "organisatieIdType" : "string",
    "timestamp" : "string",
    "abonnementId" : "string",
    "eventType": "string",
    "recordId": "aaabbbccc"
}
```


## Dependencies
 - auth