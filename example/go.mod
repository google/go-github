module github.com/google/go-github/v57/example

go 1.17

require (
	github.com/ProtonMail/go-crypto v0.0.0-20230828082145-3c4c8a2d2371
	github.com/bradleyfalzon/ghinstallation/v2 v2.0.4
	github.com/gofri/go-github-ratelimit v1.0.3
	github.com/google/go-github/v57 v57.0.0
	golang.org/x/crypto v0.17.0
	golang.org/x/term v0.15.0
	google.golang.org/appengine v1.6.7
)

require (
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-github/v41 v41.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v57 => ../
