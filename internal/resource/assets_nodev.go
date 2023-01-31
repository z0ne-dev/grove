//go:build !dev

//go:generate go run ../../scripts/generate_assets.go ../../build/resources

package resource

import (
	"net/http"
)

var All http.FileSystem = vfs
