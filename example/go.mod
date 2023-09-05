module github.com/google/go-github/v55/example

go 1.17

require (
	github.com/bradleyfalzon/ghinstallation/v2 v2.0.4
	github.com/gofri/go-github-ratelimit v1.0.3
	github.com/google/go-github/v55 v55.0.0
	golang.org/x/crypto v0.12.0
	golang.org/x/term v0.11.0
	google.golang.org/appengine v1.6.7
)

require (
	github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-github/v41 v41.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v55 => ../
