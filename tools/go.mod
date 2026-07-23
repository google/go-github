module tools

go 1.25.0

require (
	github.com/alecthomas/kong v1.15.0
	github.com/getkin/kin-openapi v0.142.0
	github.com/google/go-cmp v0.7.0
	github.com/google/go-github/v89 v89.0.0
	go.yaml.in/yaml/v3 v3.0.4
	golang.org/x/sync v0.22.0
)

require (
	github.com/go-openapi/jsonpointer v0.22.5 // indirect
	github.com/go-openapi/swag/jsonname v0.25.5 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/oasdiff/yaml v0.1.1 // indirect
	github.com/oasdiff/yaml3 v0.0.14 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	golang.org/x/text v0.14.0 // indirect
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v89 => ../
