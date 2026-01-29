module github.com/google/go-github/v81/otel/example

go 1.24.0

require (
	github.com/google/go-github/v81 v81.0.0
	github.com/google/go-github/v81/otel v0.0.0
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.24.0
	go.opentelemetry.io/otel/sdk v1.24.0
)

replace github.com/google/go-github/v81 => ../../
replace github.com/google/go-github/v81/otel => ../
