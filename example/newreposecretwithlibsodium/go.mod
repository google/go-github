module newreposecretwithlibsodium

go 1.15

require (
	github.com/GoKillers/libsodium-go v0.0.0-20171022220152-dd733721c3cb
	github.com/google/go-github/v48 v48.2.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v48 => ../..
