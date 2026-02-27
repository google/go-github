module tools

go 1.25.0

require (
	github.com/alecthomas/kong v1.14.0
	github.com/getkin/kin-openapi v0.133.0
	github.com/google/go-cmp v0.7.0
	github.com/google/go-github/v84 v84.0.0
	go.yaml.in/yaml/v3 v3.0.4
	golang.org/x/sync v0.19.0
)

require (
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-test/deep v1.1.1 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/oasdiff/yaml v0.0.0-20250309154309-f31be36b4037 // indirect
	github.com/oasdiff/yaml3 v0.0.0-20250309153720-d2182401db90 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/woodsbury/decimal128 v1.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Use version at HEAD, not the latest published.
replace github.com/google/go-github/v84 => ../
