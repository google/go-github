package github
// Provides a TokenSource that uses the basic auth mechanism to get tokens
//
//	https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
//
// Sample Usage:
//
//	oa2 := oauth2.Config{ ... }
//
//	tc, err := github.NewBasicAuthClient(oa2, "zup", "nuch", "my oauth", []string{"public_repo"})
//
//	client := github.NewClient(tc)

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Creds provides the authorization data needed to do a Basic Auth call to get a Token
// and implements the TokenSource interface
//
//	Username is the user's Basic Auth username
//	Password is the corresponding password
//	PostBodyReader is a reader that will produce the desired body for the Post for the auth call
//
// api doc showing the call is here:
//
//	https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
//
type Creds struct {
	Username       string
	Password       string
	PostBodyReader io.Reader
	Config         *oauth2.Config
}

// Token set username/password and postbody and do basic auth
func (ba Creds) Token() (tk *oauth2.Token, err error) {

	return TokenSourceBasicAuth(ba)

}

// BasicAuthRequestBody is the struct for generating the body for the authentication POST
//
// its used to generate the postBodyReader that goes into the Creds with the username/password
type BasicAuthRequestBody struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Note         string   `json:"note"`
	Scopes       []string `json:"scopes"`
}

// tokenBasicAuthJSON structure for deserializing the basic auth json response
type tokenBasicAuthJSON struct {
	Token       string `json:"token"`
	Url         string `json:"url"`
	AppRegistry struct {
		Url      string `json:"url"`
		ClientID string `json:"client_id"`
	} `json:"app_registry"`
	Note         string `json:"note"`
	RefreshToken string `json:"refresh_token"`
}

// NewBasicAuthClient returns an http client that will do the Basic Auth call to github to get a token
// via the user creds specified
//
// The returned http.Client can be passed to github.NewClient
//
// Sample Usage:
//
//	oa2 := oauth2.Config{ ... }
//
//	tc, err := github.NewBasicAuthClient(oa2, "zup", "nuch", "my oauth", []string{"public_repo"})
//
//	client := github.NewClient(tc)
//
// the github api docs are here:
//
//	https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
//
func NewBasicAuthClient(oa2 oauth2.Config, username, password, note string, repos []string) (tc *http.Client, err error) {

	postBody := BasicAuthRequestBody{
		oa2.ClientID,
		oa2.ClientSecret,
		note,
		repos,
	}

	pb, err := json.Marshal(postBody)
	if err != nil {
		return
	}

	creds := Creds{username, password, bytes.NewReader(pb), &oa2}

	tc = oauth2.NewClient(context.Background(), creds)

	return

}

// TokenSourceBasicAuth - use Basic Auth (username/password) to get a token
func TokenSourceBasicAuth(creds Creds) (*oauth2.Token, error) {

	c := creds.Config

	tk, err := retrieveTokenBasicAuth(c.Endpoint.TokenURL, creds.Username, creds.Password, creds.PostBodyReader)
	if err != nil {
		return nil, err
	}

	if tk == nil {
		return &oauth2.Token{}, fmt.Errorf("got an empty token")
	}

	return &oauth2.Token{
		AccessToken:  tk.AccessToken,
		TokenType:    tk.TokenType,
		RefreshToken: tk.RefreshToken,
		Expiry:       tk.Expiry,
	}, nil

}

// retrieveTokenBasicAuth - Do Basic auth via username/password not client_id/client_secret
//
// POST body for the auth call comes from the caller as a Reader
func retrieveTokenBasicAuth(TokenURL, Username, Password string, postBodyReader io.Reader) (*oauth2.Token, error) {

	hc := http.Client{}

	req, err := http.NewRequest("POST", TokenURL, postBodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(Username, Password)

	r, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("github: cannot fetch token: %v", err)
	}
	if code := r.StatusCode; code < 200 || code > 299 {
		return nil, fmt.Errorf("github: cannot fetch token: %v\nResponse: %s", r.Status, body)
	}

	var token *oauth2.Token
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		token = &oauth2.Token{
			AccessToken:  vals.Get("access_token"),
			TokenType:    vals.Get("token_type"),
			RefreshToken: vals.Get("refresh_token"),
		}
		e := vals.Get("expires_in")
		expires, _ := strconv.Atoi(e)
		if expires != 0 {
			token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
		}
	default:
		var tj tokenBasicAuthJSON
		if err = json.Unmarshal(body, &tj); err != nil {
			return nil, err
		}
		token = &oauth2.Token{
			AccessToken:  tj.Token,
			TokenType:    "oauth2",
			RefreshToken: tj.Token,
			Expiry:       time.Now().Add(time.Duration(99) * time.Minute),
		}
	}

	return token, nil
}
