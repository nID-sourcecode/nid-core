// Package graphql resolvers for the graphql server.
package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	"golang.org/x/crypto/bcrypt"
	"lab.weave.nl/weave/generator/gen/auth"
	jwtutils "lab.weave.nl/weave/generator/jwtutils"
	generr "lab.weave.nl/weave/generator/pkg/errors"
)

var (
	// ErrWrongCredentials graphql invalid credentials
	ErrWrongCredentials = generr.NewGraphQLError("wrong credentials", false, "GEN-AUTH-001")
	// ErrTokenExpired graphql token expired
	ErrTokenExpired = generr.NewGraphQLError("token expired", false, "GEN-AUTH-002")
	// ErrTokenNotFound graphql token not found
	ErrTokenNotFound = generr.NewGraphQLError("token not found", false, "GEN-AUTH-003")
	// ErrCouldNotParseToken graphql failed to parse jwt token
	ErrCouldNotParseToken = generr.NewGraphQLError("failed to parse jwt token", false, "GEN-AUTH-004")
	// ErrInvalidToken graphql invalid jwt token
	ErrInvalidToken = generr.NewGraphQLError("invalid jwt token", false, "GEN-AUTH-005")
	// ErrInvalidClaims graphql invalid claims
	ErrInvalidClaims = generr.NewGraphQLError("invalid claims", false, "GEN-AUTH-006")
	// ErrFailedToRetrieveUser graphql failed to retrieve user
	ErrFailedToRetrieveUser = generr.NewGraphQLError("failed to retrieve user", true, "GEN-AUTH-007")
	// ErrGettingTokenFromDB graphql failed to get token from db
	ErrGettingTokenFromDB = generr.NewGraphQLError("getting token from db", true, "GEN-AUTH-008")
)

func (r *mutationResolver) RefreshJwt(_ context.Context, token string) (string, error) {
	err := validateToken(r.DB, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) CreateJwt(_ context.Context, email string, password string) (string, error) {
	user, err := checkUser(r.DB, email, password)
	if err != nil {
		return "", err
	}

	return createJWTToken(r.DB, user)
}

func checkUser(db *gorm.DB, email string, password string) (*models.User, error) {
	// Check input
	if email == "" || password == "" {
		return nil, ErrWrongCredentials
	}

	// Find user
	var user models.User
	email = strings.ToLower(email)
	err := db.Model(user).Where("lower(email) = ? AND password IS NOT NULL AND password != ''", email).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWrongCredentials
		}
		return nil, fmt.Errorf("%w \"%s\"", ErrFailedToRetrieveUser, email)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrWrongCredentials
	}

	return &user, nil
}

func validateToken(db *gorm.DB, token string) error {
	// Parse jwt token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return jwtutils.PublicKey(), nil })
	if err != nil {
		return fmt.Errorf("%w \"%s\"", ErrCouldNotParseToken, token)
	}
	// Valid token check
	if !parsedToken.Valid {
		return fmt.Errorf("%w \"%s\"", ErrInvalidToken, token)
	}
	// Jwt token claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("%w \"%v\" for jwt token \"%s\"", ErrInvalidClaims, parsedToken.Claims, token)
	}

	var jwtFromDB auth.JWT
	// Check if jwt still exists in db
	err = db.Model(&jwtFromDB).Where("id = ?", claims["jti"]).Find(&jwtFromDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w \"%s\"", ErrTokenNotFound, claims["jti"])
		}
		return fmt.Errorf("%w \"%s\": %v", ErrGettingTokenFromDB, claims["jti"], err.Error())
	}

	// Is the found jwt not expired
	var expTime time.Time
	if exp, ok := claims["exp"].(float64); ok {
		expTime = time.Unix(int64(exp), 0)
	}
	if expTime.Before(time.Now()) {
		return fmt.Errorf("%w \"%s\"", ErrTokenExpired, claims["jti"])
	}

	return nil
}

func createJWTToken(db *gorm.DB, user *models.User) (string, error) {
	var scopes []string
	err := json.Unmarshal(user.Scopes.RawMessage, &scopes)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodRS512)

	uniqueID := uuid.Must(uuid.NewV4())
	now := time.Now()
	iat := now.Unix()
	exp := now.Add(time.Hour * 24 * 7 * 4) // in 4 weeks
	nbf := now.Add(-2 * time.Minute)       // 2 minutes ago
	token.Claims = jwt.MapClaims{
		"jti":    uniqueID,   // a unique identifier for the token
		"exp":    exp.Unix(), // time when the token will expire (x minutes from now)
		"iat":    iat,        // when the token was issued/created (now)
		"nbf":    nbf.Unix(), // time before which the token is not yet valid (2 minutes ago)
		"sub":    user.ID,    // the subject/principal is whom the token is about
		"user":   user.ID,    // user id
		"scopes": scopes,     // token scope - not a standard claim
	}

	signedToken, err := token.SignedString(jwtutils.PrivateKey())
	if err != nil {
		return "", generr.WrapAsInternal(err, "signing token")
	}

	// Create JWTModel for DB
	var jwtmodel auth.JWT
	jwtmodel.ID = uniqueID
	jwtmodel.UserID = fmt.Sprintf("%v", user.ID)
	jwtmodel.Expiration = exp

	if err := db.FirstOrCreate(&jwtmodel).Error; err != nil {
		return "", generr.WrapAsInternal(err, "creating jwt")
	}

	return signedToken, nil
}
