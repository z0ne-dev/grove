package resource

import (
	"net/http"
	"path/filepath"
	"strings"
)

type PrefixFileSystem struct {
	fs     http.FileSystem
	prefix string
}

func (fs *PrefixFileSystem) Open(name string) (http.File, error) {
	prefixedName := filepath.Join(fs.prefix, name)

	return fs.fs.Open(prefixedName)
}

func NewPrefixFileSystem(fs http.FileSystem, prefix string) http.FileSystem {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	return &PrefixFileSystem{
		fs:     fs,
		prefix: prefix,
	}
}
