package resource

import (
	"net/http"
)

var (
	Assets     http.FileSystem
	Templates  http.FileSystem
	Migrations http.FileSystem
)

func init() {
	Assets = NewPrefixFileSystem(All, "/assets/")
	Templates = NewPrefixFileSystem(All, "/templates/")
	Migrations = NewPrefixFileSystem(All, "/migrations/")
}
