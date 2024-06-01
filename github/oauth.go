package github

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

const (
	// DefaultApplicationTokenExpiration is the default expiration time for the GitHub App token.
	// The expiration time of the JWT, after which it can't be used to request an installation token.
	// The time must be no more than 10 minutes into the future.
	DefaultApplicationTokenExpiration = 10 * time.Minute

	bearerTokenType = "Bearer"
)

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
		TokenType:   bearerTokenType,
		Expiry:      expiresAt,
	}, nil
}

// InstallationTokenSourceOpt is a functional option for InstallationTokenSource.
type InstallationTokenSourceOpt func(*installationTokenSource)

// WithInstallationTokenOptions sets the options for the GitHub App installation token.
func WithInstallationTokenOptions(opts *InstallationTokenOptions) InstallationTokenSourceOpt {
	return func(i *installationTokenSource) {
		i.opts = opts
	}
}

// InstallationTokenSource represents a GitHub App installation token source.
type installationTokenSource struct {
	id   int64
	src  oauth2.TokenSource
	apps *AppsService
	opts *InstallationTokenOptions
}

// NewInstallationTokenSource creates a new GitHub App installation token source.
func NewInstallationTokenSource(id int64, src oauth2.TokenSource, opts ...InstallationTokenSourceOpt) oauth2.TokenSource {
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: src,
		},
	}

	i := &installationTokenSource{
		id:   id,
		src:  src,
		apps: NewClient(client).Apps,
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

// Token creates a new GitHub App installation token.
// The token is used to authenticate as a GitHub App installation.
func (t *installationTokenSource) Token() (*oauth2.Token, error) {
	ctx := context.Background()

	token, _, err := t.apps.CreateInstallationToken(ctx, t.id, t.opts)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: token.GetToken(),
		TokenType:   bearerTokenType,
		Expiry:      token.GetExpiresAt().Time,
	}, nil
}
