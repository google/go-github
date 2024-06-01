package github

import (
	"crypto/rsa"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

// DefaultApplicationTokenExpiration is the default expiration time for the GitHub App token.
// The expiration time of the JWT, after which it can't be used to request an installation token.
// The time must be no more than 10 minutes into the future.
const DefaultApplicationTokenExpiration = 10 * time.Minute

// applicationTokenSource represents a GitHub App token.
// https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app
type applicationTokenSource struct {
	id         string
	privateKey *rsa.PrivateKey
	expiration time.Duration
}

// ApplicationTokenOpt is a functional option for ApplicationTokenSource.
type ApplicationTokenOpt func(*applicationTokenSource)

// WithApplicationTokenExpiration sets the expiration for the GitHub App token.
// The expiration time of the JWT, after which it can't be used to request an installation token.
// The time must be no more than 10 minutes into the future.
func WithApplicationTokenExpiration(expiration time.Duration) ApplicationTokenOpt {
	return func(a *applicationTokenSource) {
		// The expiration time must be no more than 10 minutes into the future.
		// Also, the expiration time must be greater than 0.
		if expiration > DefaultApplicationTokenExpiration || expiration <= 0 {
			expiration = DefaultApplicationTokenExpiration
		}
		a.expiration = expiration
	}
}

// NewApplicationTokenSource creates a new GitHub App token source.
// An application token is used to authenticate as a GitHub App.
func NewApplicationTokenSource(id string, privateKey []byte, opts ...ApplicationTokenOpt) (oauth2.TokenSource, error) {
	if id == "" {
		return nil, errors.New("applicationID is required")
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	t := &applicationTokenSource{
		id:         id,
		privateKey: privKey,
		expiration: DefaultApplicationTokenExpiration,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t, nil
}

// Token creates a new GitHub App token.
// The token is used to authenticate as a GitHub App.
// Each token is valid for 10 minutes.
func (t *applicationTokenSource) Token() (*oauth2.Token, error) {
	// To protect against clock drift, set this 60 seconds in the past.
	now := time.Now().Add(-60 * time.Second)
	expiresAt := now.Add(t.expiration)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    t.id,
	})

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		Expiry:      expiresAt,
	}, nil
}
