// assets_nodev.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

//go:build !dev

//go:generate sh ../scripts/generate.sh

package resource

import "embed"

var (
	//go:embed frontend/build
	Assets embed.FS
	//go:embed templates
	Templates embed.FS
	//go:embed migrations
	Migrations embed.FS
)
