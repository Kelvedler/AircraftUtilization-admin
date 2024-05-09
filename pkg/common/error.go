package common

import (
	"html/template"
	"net/http"
)

func ErrorResp(w http.ResponseWriter) {
	errMsg := "Something Went Wrong"
	data := make(map[string]string)
	data["Message"] = errMsg
	w.Header().Set("HX-Retarget", "#content")
	tmpl := template.Must(template.ParseFiles("template/base.html")).Lookup("error-page")
	tmpl.Execute(w, data)
}
