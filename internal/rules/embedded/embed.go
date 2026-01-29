package embedded

import _ "embed"

//go:embed core.yaml
var CoreYAML []byte

//go:embed credentials.yaml
var CredentialsYAML []byte

//go:embed security.yaml
var SecurityYAML []byte

//go:embed packages.yaml
var PackagesYAML []byte

//go:embed configuration.yaml
var ConfigurationYAML []byte
