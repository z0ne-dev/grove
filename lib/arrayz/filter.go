// filter.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package arrayz

func Filter[T any](a []T, fn func(T) bool) []T {
	var r []T
	for _, v := range a {
		if fn(v) {
			r = append(r, v)
		}
	}
	return r
}
