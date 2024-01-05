package config

import _ "embed"

//go:embed services.json
var ServicesConfigJSON []byte
