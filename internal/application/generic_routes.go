package application

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi"
	"net/http"
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
