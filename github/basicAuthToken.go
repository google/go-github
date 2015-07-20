package github

import (
	"fmt"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	// "strings"
	"time"

	// "golang.org/x/net/context"
	"golang.org/x/oauth2"
)
// basicAuthRequestBody is the struct for generating the body for the post
//
// its used to generate the postBodyReader that goes into the oauth2Creds with the username/password
type BasicAuthRequestBody struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Note         string   `json:"note"`
	Scopes       []string `json:"scopes"`
}

// tokenBasicAuthJSON basic auth json response
//
// it has no RefreshToken, Expires*
// leave them in the struct so it won't break the expiry call and maybe other calls
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

// BasicAuth is a struct which is used to orchestrate the Token call
//
// Token calls take no arguments, but the creds are in the context
// and the oauth2 config is needed too
// type BasicAuth struct {
// 	context.Context
// 	oauth2.Config
// }

// Creds provides the authorization data needed to do a Basic Auth call to get a token
//
// 	Username is Basic Auth username
// 	Password is the corresponding password
// 	PostBodyReader is a reader that will produce the desired body for the Post for the auth call
type Creds struct {
  Username string
  Password string
  PostBodyReader io.Reader
  Config *oauth2.Config
}

// CredsGetter is a type that will get creds from a context for one of these calls
// type CredsGetter func(context.Context) (creds, bool)

// NewBasicAuthClient generates an http client that will do the Basic Auth call to github to get a token
//
// The returned client can be passed to github.NewClient
//
// the api docs are here https://developer.github.com/v3/oauth_authorizations/#create-a-new-authorization
//
// func NewBasicAuthClient(oa2 oauth2.Config, username, password, note string, repos []string) (tc *http.Client, err error) {

// 	postBody := basicAuthRequestBody{
// 		oa2.ClientID,
// 		oa2.ClientSecret,
// 		note,
// 		repos,
// 	}

// 	pb, err := json.Marshal(postBody)
// 	if err != nil {
// 		return
// 	}

// 	creds := creds{username, password, bytes.NewReader(pb)}

// 	ctx := NewContext(context.Background(), creds)

// 	tc = oauth2.NewClient(ctx, BasicAuth{ctx, oa2})

// 	return

// }

//set username/password and postbody in the context
//
//needs to get client from context so this ting needs to get context somehow
func (ba Creds) Token() (tk *oauth2.Token, err error) {

	return TokenSourceBasicAuth(ba)

}

// GetTokenBasicAuth - use Basic Auth (username/password) to get a token
//
// Get a token with basic auth as in:
//
// https://developer.github.com/enterprise/2.1/v3/oauth_authorizations/#get-or-create-an-authorization-for-a-specific-app
//
func TokenSourceBasicAuth(creds Creds) (*oauth2.Token, error) {

	c := creds.Config
	// creds, ok := FromContext(ctx)

	// if !ok {
	// 	return nil, fmt.Errorf("creds for basic auth not found")
	// }

	tk, err := retrieveTokenBasicAuth(c.Endpoint.TokenURL, creds.Username, creds.Password, creds.PostBodyReader)
	if err != nil {
		return nil, err
	}

	// return tokenFromInternal(tk), nil
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

// Do Basic auth via username/password not client_id/client_secret
//
// POST body for the auth call comes from the caller as a Reader
func retrieveTokenBasicAuth(TokenURL, Username, Password string, postBodyReader io.Reader) (*oauth2.Token, error) {

	// hc, err := ContextClient(ctx)
	// if err != nil {
	// 	return nil, err
	// }

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
		if e == "" {
			// TODO(jbd): Facebook's OAuth2 implementation is broken and
			// returns expires_in field in expires. Remove the fallback to expires,
			// when Facebook fixes their implementation.
			e = vals.Get("expires")
		}
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

//typesafe context acccessors
// type key int

// var CredsKey key = 0

// func NewContext(ctx context.Context, ba Creds) context.Context {
// 	return context.WithValue(ctx, CredsKey, ba)
// }

// func FromContext(ctx context.Context) (Creds, bool) {
// 	ba, ok := ctx.Value(CredsKey).(Creds)
// 	return ba, ok
// }
