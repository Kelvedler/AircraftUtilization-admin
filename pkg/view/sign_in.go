package view

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/auth"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/common"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/crypto"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/db"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/middleware"
)

func sanitizeAdmin(rc *middleware.RequestContext, a *db.Admin) {
	sanitizer := rc.Sanitize
	a.Name = sanitizer.Sanitize(a.Name)
	a.Password = sanitizer.Sanitize(a.Password)
}

func signInApi(
	rc *middleware.RequestContext,
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	tmpl := template.Must(template.ParseFiles("template/index-component.html")).
		Lookup("sign-in-form")
	errMap := make(map[string]string)
	errMap["InputErr"] = "Invalid credentials"
	var admin db.Admin
	err := common.BindJSON(r, &admin)
	if err != nil {
		rc.Logger.Error(err.Error())
		common.ErrorResp(w)
		return
	}
	sanitizeAdmin(rc, &admin)
	err = rc.Validate.StructPartial(admin, "Name", "Password")
	if err != nil {
		rc.Logger.Info(err.Error())
		tmpl.Execute(w, errMap)
		return
	}
	rawPass := admin.Password
	errs := db.PerformBatch(r.Context(), rc.DbPool, []db.BatchSet{admin.GetByName})
	adminErr := errs[0]
	if adminErr != nil {
		rc.Logger.Info(fmt.Sprintf("Admin %s not found", admin.Name))
		tmpl.Execute(w, errMap)
		return
	}
	correctPass, err := crypto.CompareKeys(rawPass, admin.Password)
	if !correctPass {
		rc.Logger.Info("Incorrect password")
		tmpl.Execute(w, errMap)
		return
	}
	err = auth.SetNewTokenCookie(w, admin)
	if err != nil {
		rc.Logger.Error(err.Error())
		common.ErrorResp(w)
		return
	}
	w.Header().Set("HX-Redirect", "/api-users")
}
