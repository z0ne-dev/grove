// generic_routes.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package application

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi"
)

type GenericRoutes struct {
	set *jet.Set
}

func NewGenericRoutes(set *jet.Set) *GenericRoutes {
	return &GenericRoutes{
		set: set,
	}
}

func (route *GenericRoutes) Routes(r chi.Router) {
	r.Get("/", route.GetHome)
}

func (route *GenericRoutes) GetHome(writer http.ResponseWriter, _ *http.Request) {
	template, err := route.set.GetTemplate("_index.jet")
	if err != nil {
		println(err.Error())
		return
	}

	err = template.Execute(writer, nil, nil)
	if err != nil {
		println(err.Error())
		return
	}
}
