// sort.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package arrayz

import (
	"sort"
)

var _ sort.Interface = (*sortInterface[int])(nil)

type sortInterface[T any] struct {
	data    []T
	compare func(T, T) bool
}

func (s *sortInterface[T]) Len() int {
	return len(s.data)
}

func (s *sortInterface[T]) Less(i, j int) bool {
	return s.compare(s.data[i], s.data[j])
}

func (s *sortInterface[T]) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func Sort[T any](a []T, fn func(T, T) bool) []T {
	if len(a) < 2 {
		return a
	}

	sorter := &sortInterface[T]{
		data:    a,
		compare: fn,
	}
	sort.Sort(sorter)

	return sorter.data
}
