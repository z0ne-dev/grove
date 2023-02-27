// assets_dev.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package resource

import (
	"os"
)

var (
	Assets     = os.DirFS("frontend/build")
	Templates  = os.DirFS("templates")
	Migrations = os.DirFS("migrations")
)
