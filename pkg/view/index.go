package view

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/middleware"
)

func index(
	_ *middleware.RequestContext,
	w http.ResponseWriter,
	_ *http.Request,
	_ httprouter.Params,
) {
	tmpl := template.Must(
		template.ParseFiles(
			"template/index.html",
			"template/index-component.html",
			"template/base.html",
		),
	)
	tmpl.Execute(w, nil)
}
