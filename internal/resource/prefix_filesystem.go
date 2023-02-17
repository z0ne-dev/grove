// prefix_filesystem.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

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
