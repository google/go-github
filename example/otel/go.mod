module github.com/google/go-github/v81/example/otel

go 1.24.0

require (
	github.com/google/go-github/v81 v81.0.0
	github.com/google/go-github/v81/otel v0.0.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.27.0
	go.opentelemetry.io/otel/sdk v1.27.0
)

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	go.opentelemetry.io/otel v1.27.0 // indirect
	go.opentelemetry.io/otel/metric v1.27.0 // indirect
	go.opentelemetry.io/otel/trace v1.27.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

replace github.com/google/go-github/v81 => ../../

replace github.com/google/go-github/v81/otel => ../../otel
