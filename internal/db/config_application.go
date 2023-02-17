// config_application.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package db

type ConfigPayloadApplication struct {
	Name string `json:"name" default:"Grove"`
}
