// mapstructure_json_unmarshaller_hook_func.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package util

import (
	"encoding/json"
	"reflect"
	"regexp"

	"github.com/mitchellh/mapstructure"
)

// reMatchJsonTypes is a very simple regular expression that is supposed to match valid json.
var reMatchJsonTypes = regexp.MustCompile(`^([\[{"]|\d+$|null$|true$|false$)`)

// MapStructureJSONUnmarshallerHookFunc is a mapstructure.DecodeHookFuncType that unmarshals using json.Unmarshaler.
func MapStructureJSONUnmarshallerHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		result := reflect.New(t).Interface()
		unmarshaller, ok := result.(json.Unmarshaler)
		if !ok {
			return data, nil
		}

		str, ok := data.(string)
		if !ok {
			return data, nil
		}

		if str == "" {
			str = "null"
		} else if !reMatchJsonTypes.MatchString(str) {
			str = `"` + str + `"`
		}

		if err := unmarshaller.UnmarshalJSON([]byte(str)); err != nil {
			return nil, err
		}
		return result, nil
	}
}
