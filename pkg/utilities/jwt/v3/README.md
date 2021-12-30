# JWT

A weave implementation / extension of the [jwt-go]((https://github.com/dgrijalva/jwt-go) package.

## Installation

Install using: `go get -u lab.weave.nl/weave/utilities/jwt/v3@v3.0.0`

## Usage

1. You first need to initiate the package by specifying a public and private key 

```go
client := jwt.NewJWTClient(priv, pub)
```
By default de DefaultOpts() is called which will create the default ClientOps containing information about the signing method and the JWT header configuration

2. Specify the claims either by using the DefaultClaims or a custom implementation of the Claims interface:

```go
defaultClaims := jwt.NewDefaultClaims()
```

3. Create the token:

```go
token, err := jwt.SignToken(defaultClaims)
if err != nil {
        return err

```

4. Verify that the token is correct:

```go
defaultClaims := &DefaultClaims{}
err := jwt.ValidateAndParse(token, defaultClaims)
if err != nil {
        return err
}
```

The ValidateAndParse will validate the token and parse it into de defaultClaims. Therefore the type assertion is automatically done.

## Custom Claims

The DefaultClaims conform to the RFC specification. If you want to add additional properties or create your custom claims you can do that as follows:

### Extend Default Claims

1. Create your custom struct and extend the DefaultClaims

```go
type CustomClaims struct {
        *DefaultClaims
        AddionalProperty `json:"additional_property"`
}
```

2. Since the DefaultClaims conforms to the Claims interface you can already use the CustomClaims

3. If you want to add a additional validator you can do:

```go
func (c *CustomClaims) Valid() error {
    err := c.defaultClaims.Valid()
    if err != nil {
        return nil
    }

    // CUSTOM VALIDATION
    ...
}
```

### Custom struct conform to the Claims interface

When the properties of the defaultClaims are not needed you can also create your custom claims. The custom claims struct needs to implement the Valid() and ParseToken() methods to conform to the Claims interface.
