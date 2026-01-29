module github.com/google/go-github/v82/example/otel

go 1.24.0

require (
	github.com/google/go-github/v82 v82.0.0
	github.com/google/go-github/v82/otel v0.0.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.24.0
	go.opentelemetry.io/otel/sdk v1.24.0
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
)

replace github.com/google/go-github/v82 => ../../

replace github.com/google/go-github/v82/otel => ../../otel
