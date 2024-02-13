package github

// HookConfig describes metadata about a webhook configuration.
type HookConfig struct {
	// The media type used to serialize the payloads
	// Possible values are `json` and `form`, the field is not specified the default is `form`
	ContentType *string `json:"content_type,omitempty"`
	// The possible values are 0 and 1.
	// Setting it to 1 will allow skip certificate verification for the host,
	// potentially exposing to MitM attacks: https://en.wikipedia.org/wiki/Man-in-the-middle_attack
	InsecureSSL *string `json:"insecure_ssl,omitempty"`
	URL         *string `json:"url,omitempty"`

	// Secret is returned obfuscated by GitHub, but it can be set for outgoing requests.
	Secret *string `json:"secret,omitempty"`
}
