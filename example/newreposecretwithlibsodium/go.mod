module newreposecretwithlibsodium

go 1.22.10

require (
	github.com/GoKillers/libsodium-go v0.0.0-20171022220152-dd733721c3cb
	github.com/google/go-github/v68 v68.0.0
)

require github.com/google/go-querystring v1.1.0 // indirect

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v68 => ../..
