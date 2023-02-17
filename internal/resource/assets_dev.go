// assets_dev.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package resource

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/httpfs/union"
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
