module github.com/google/go-github/v85/tools/check-structfield-settings

go 1.25.0

require (
	github.com/golangci/plugin-module-register v0.1.2
	github.com/google/go-github/v85/tools/structfield v0.0.0
	go.yaml.in/yaml/v3 v3.0.4
	golang.org/x/tools v0.44.0
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v85/tools/structfield v0.0.0 => ../structfield
