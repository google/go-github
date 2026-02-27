module newreposecretwithlibsodium

go 1.25.0

require (
	github.com/GoKillers/libsodium-go v0.0.0-20171022220152-dd733721c3cb
	github.com/google/go-github/v84 v84.0.0
)

require github.com/google/go-querystring v1.2.0 // indirect

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v84 => ../..
