package resource

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
	"os"
	"path/filepath"
)

var All http.FileSystem

func init() {
	wd, _ := os.Getwd()
	All = union.New(map[string]http.FileSystem{
		"/templates":  http.Dir(filepath.Join(wd, "templates")),
		"/assets":     http.Dir(filepath.Join(wd, "frontend", "build")),
		"/migrations": http.Dir(filepath.Join(wd, "migrations")),
	})
}
