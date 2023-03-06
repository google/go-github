module github.com/google/go-github/v50/example

go 1.17

require (
	github.com/bradleyfalzon/ghinstallation/v2 v2.0.4
	github.com/gofri/go-github-ratelimit v1.0.1
	github.com/google/go-github/v50 v50.0.0
	golang.org/x/crypto v0.1.0
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	google.golang.org/appengine v1.6.7
)

require (
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/go-github/v41 v41.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.1.0 // indirect
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v50 => ../
