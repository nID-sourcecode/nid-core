# Retrieving data
## Prerequisites
- GraphQLServiceURL (example: https://[domain])
- AccessToken (received after going trough a successful auth flow)

**Note: more information about the auth flow can be found on the auth service flow documentation page.**

## Fetching data
In this example we are going to fetch users bank account and contact details. Because the service we use in this example is a GraphQL service we send all requests to a single endpoint. In our case this endpoint is the url of the GraphQLService followed by `/gql` and is reachable via a `POST` or `GET` request.

### Querying user data
Before querying the data we need to specify a query otherwise the GraphQL service doesn't know what data we are asking for. <Br>
**Bank account query:**  <br>
```graphql
query {
    users(filter: {
        pseudonym: {
            eq: "$$nid:subject$$"
        }
    }){
        bankAccounts {
            savingsAccounts {
                name
                amount
            }
            accountNumber
            amount
        }
    }
}
``` 
**Contact details query:**  <br>
```graphql
query {
    users(filter: {
        pseudonym: {
            eq: "$$nid:subject$$"
        }
    }){
        contactDetails {
            address {
                houseNumber
                houseNumberAddon
                postalCode
            }
            phone
        }
    }
}
``` 
The queries mentioned here are default GraphQL queries expected for the filter value `"$$nid:subject$$"`. This filter value is necessary, without it we can't fetch user data. You don't have to explicitly set this value, the system will replace this based on the values in the provided JWT token.
### Authorization
Sending a query to the databron won't work without an authorization header. Every databron request expects an `Authorization` header with the access token as value. For example: `Authorization: Bearer [AccessToken]` Make sure the `Bearer` is used as prefix. With the accessToken in place the GraphQL filter mentioned before will be correctly translated to the corresponding user identifier.
### Examples
**1. POST request to fetch users contact details**
```
> POST /gql HTTP/1.1
> Host: [domain]
> Content-Type: application/json
> Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJodHRwOi8vb3RoZXItZGF0YWJyb24ubmlkL2dxbCIsImNsaWVudF9pZCI6IjZiYTdiODEwLTlkYWQtMTFkMS04MGI0LTAwYzA0ZmQ0MzBjOCIsImV4cCI6MTYwMzkwMDAwNiwiaWF0IjoxNjAzODEzNjA2LCJpc3MiOiJhdXRoLm5pZCIsImp0aSI6Ijc2NTE1M2I1LWEyMjAtNGY1NC1hYmE4LTU1ZmJiYjdlMmUzYSIsIm5iZiI6MTYwMzgxMzQ4NiwicXVlcmllcyI6eyJkYXRhYnJvbjphZGRyZXNzZXMiOnsiQyI6eyJmIjpbInBob25lIl19LCJVIjp7ImYiOltdLCJtIjp7ImNvbnRhY3REZXRhaWxzIjoiI0MifSwicCI6eyJmaWx0ZXIiOnsicHNldWRvbnltIjp7ImVxIjoiJCRuaWQ6c3ViamVjdCQkIn19fX0sInIiOnsibSI6eyJ1c2VycyI6IiNVIn19fSwiZGF0YWJyb246Y29udGFjdGRldGFpbHMiOnsiQSI6eyJmIjpbImhvdXNlTnVtYmVyIl19LCJDIjp7ImYiOltdLCJtIjp7ImFkZHJlc3MiOiIjQSJ9fSwiVSI6eyJmIjpbXSwibSI6eyJjb250YWN0RGV0YWlscyI6IiNDIn0sInAiOnsiZmlsdGVyIjp7InBzZXVkb255bSI6eyJlcSI6IiQkbmlkOnN1YmplY3QkJCJ9fX19LCJyIjp7Im0iOnsidXNlcnMiOiIjVSJ9fX19LCJzdWIiOiJXUnY3dE1tQUp2NzBFeGNMU1N1NjgzL2ZLRDEwR0lhWXdvY0xZREJvcG81TVJVaHpXdWxwWmtZOENLckVHMEVibUFoeHZxN1dWU1ZidGhkckRtdmdwWXphSjBMRWJvZ2FHNWJldTlJdVNQbmlmc0lvUDF2ems2Sk9KaVNNRnZvaHlYKzRyeVZLam1URFhHY3pCWTRqVFV0dXlYZ0ROQ2xuT2EwdkRldnV1cHVXTHhyNzUxZDd4bDFWRVUvam94eU1VTVVkVU9KczBSSytXNURKVU1DaGpQcXVaQXVFNUQzRjdPMkNqb3JLdE9lT0hUZjJjVFVDTHp1cnZ2VUVvajB3VksvSENkdjJjR3lobDhaajZjRysvMWE3L0tNYUNNaVR2cHZkcFg3eXA5TGdnQ1JVZnNxQ2RiK2Jrb1k3eGwxM1g0aUZVUGJlTVBTUjRBVUZYRURDWnc9PSIsInN1YmplY3RzIjp7Im5pZCI6IlNSOUhBNlo3cFNwUGNZM05yK3VzZXVId0YzREsrUXZVMG1aNUxVZ0NjcytxQW5BQ0pQYldtRjNMMEpMYm02eFFCbHdnbXdCeW8xZWt1MlFlVk9pb3k5bmVMT3hQR21aRnpwZDRlVkZJakJJMkhUZ0tpUWdLbEk1ZTNXeFVOcWFPVC9vL1FkU1Y3MUFQOFowZEVvZjNmcEU2VituVysvMWQ0ZHJuNldWS0NiSGJqT3ZNekIzVkNuY2dYWlExK3YyRzFBMnRLeDk4QTRTSkRrREVXajZ2TlowTHREMzlsME9xN3FaL3p1bkJGZ3dQU2FQK3NKa2xzNXpvMmZmTTk5dUdxaVlQTHZsSHVkZThjTHVMaVQra3h1Rk5ua2dwNEgvRnU0THVTK1JxdDIvRXhBSjlBc1V5elh3ZEZIbEpNSnlxZWJQblJVUXAzZ1JMR3lJMVJYckM5dz09In19.klTfEA4G8ISpJ8IQfxDpNi96lhEAo1-bvGM0A7lm6SDfQ14TU6Zh3Ytu3e6eF2QTKous1hOduXlm_q001UlnpF6tI117MKthN_ZnrTc3hSxVpNX4KN0JgKTIZE2dXUCdQa-OoGfDetapCrfR8G2Wp7Dw8m5DjSdvbtgPUXtEIzweaWiK3YpluLnwFLrDKBXpOpXplnEPrF19ZB955EU_O3X0LpTs5IYOObCYWKoaz5sJEiacH5bkFmD8amuNh6GjX6c6bR-u4gg5EreZKmRMVLdidF423HMCjwiQ-k4VPcFyRWM8vFmNNwioYivk0KbsVXh9_sTC5rPKGiZRj_HJ5A

| {"query":"query {\n    users(filter: {\n        pseudonym: {\n            eq: \"$$nid:subject$$\"\n        }\n    }){\n        contactDetails {\n            address {\n                houseNumber\n            }\n        }\n    }\n}"}

```
### Response codes and messages explained
#### Error response messages
**400 BadRequest:** Expect this response when the `Authorization` header is not send with the request, send with the request but empty or does not contain the `Bearer` prefix.
```json
{
  "errors": [
    {
      "message": "Authorization header \"\" does not adhere to the Bearer scheme"
    }
  ]
}
```
**403 Forbidden:** Expect this message when you are querying data for a user **without permission from the user**. This means the user did not gave consent to query data from the model found in the query.
```json
{
  "errors": [
    {
      "message": "query does not match scope"
    }
  ]
}
```
#### Successful response messages
**200 OK:** Queried data found.
```json
{
  "data": {
    "users": [
      {
        "contactDetails": [
          {
            "address": {
              "houseNumber": 1001
            }
          }
        ]
      }
    ]
  }
}
```
**NOTE: The response data will vary based on the query found in de request body.**<br> 
**200 OK:** Queried data not found in the response. This means the user doesn't have contact details in this case.  
```json
{
  "data": {
    "users": []
  }
}
``` 