module newreposecretwithlibsodium

go 1.15

require (
	github.com/GoKillers/libsodium-go v0.0.0-20171022220152-dd733721c3cb
	github.com/google/go-github/v56 v56.0.0
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v56 => ../..
