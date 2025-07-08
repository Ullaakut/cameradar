package dict

import (
	_ "embed"
)

//go:embed assets/credentials.json
var defaultCredentials []byte

//go:embed assets/routes
var defaultRoutes string
