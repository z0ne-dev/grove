// assets_prefixes.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

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
