// generate_assets.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/httpfs/union"

	"github.com/shurcooL/vfsgen"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	wd, err := os.Getwd()
	must(err)

	println(wd)

	assets := filepath.Join(wd, "frontend", "build")
	templates := filepath.Join(wd, "templates")
	migrations := filepath.Join(wd, "migrations")

	fs := union.New(map[string]http.FileSystem{
		"/templates":  http.Dir(templates),
		"/assets":     http.Dir(assets),
		"/migrations": http.Dir(migrations),
	})

	must(vfsgen.Generate(fs, vfsgen.Options{
		Filename:     filepath.Join(wd, "assets_generated.go"),
		PackageName:  "resource",
		VariableName: "All",
		BuildTags:    "!dev",
	}))
}
