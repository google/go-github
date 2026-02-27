module github.com/google/go-github/v84/tools/check-structfield-settings

go 1.25.0

require (
	github.com/golangci/plugin-module-register v0.1.1
	github.com/google/go-github/v84/tools/structfield v0.0.0
	go.yaml.in/yaml/v3 v3.0.4
	golang.org/x/tools v0.29.0
)

require (
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v84/tools/structfield v0.0.0 => ../structfield
