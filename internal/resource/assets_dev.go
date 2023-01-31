package resource

import (
	"net/http"
	"os"
	"path/filepath"
)

var All http.FileSystem

func init() {
	wd, _ := os.Getwd()
	All = http.Dir(filepath.Join(wd, "build", "resources"))
}
