module tools

go 1.21

toolchain go1.22.0

require (
	github.com/alecthomas/kong v0.9.0
	github.com/getkin/kin-openapi v0.126.0
	github.com/google/go-cmp v0.6.0
	github.com/google/go-github/v63 v63.0.0
	golang.org/x/sync v0.7.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/invopop/yaml v0.3.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v63 => ../
