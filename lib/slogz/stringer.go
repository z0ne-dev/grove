// stringer.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package slogz

import (
	"fmt"

	"golang.org/x/exp/slog"
)

func Stringer(key string, value any) slog.Attr {
	str := ""
	if s, ok := value.(fmt.Stringer); ok {
		str = s.String()
	} else {
		str = fmt.Sprintf("%v", value)
	}

	return slog.String(key, str)
}
