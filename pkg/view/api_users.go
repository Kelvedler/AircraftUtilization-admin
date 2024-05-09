package view

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/xsrftoken"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/common"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/crypto"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/db"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/middleware"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

type apiUsersData struct {
	ApiUsers     []db.ApiUser
	Name         string
	Key          string
	CreatedState bool
	Page         uint8
	PreviousPage string
	NextPage     string
	InputErr     string
	PostXsrf     string
}

func postApiUserXsrf(adminId uuid.UUID) string {
	return xsrftoken.Generate(
		setting.Setting.SecretKey,
		adminId.String(),
		"/api/v1/api-users",
	)
}

func (usersData *apiUsersData) populateFullPage() {
	usersToFull := int(setting.Setting.PageSize) - len(usersData.ApiUsers)
	for i := 0; i < usersToFull; i++ {
		usersData.ApiUsers = append(usersData.ApiUsers, db.ApiUser{})
	}
}

func (usersData apiUsersData) nextPageExists() bool {
	if uint8(len(usersData.ApiUsers)) > setting.Setting.PageSize {
		return true
	}
	return false
}

func sanitizeApiUser(rc *middleware.RequestContext, apiUser *db.ApiUser) {
	sanitizer := rc.Sanitize
	apiUser.Name = sanitizer.Sanitize(apiUser.Name)
}

func apiUsersApiUrl(page uint8) string {
	return fmt.Sprintf("/api/v1/api-users/?page=%d", page)
}

func apiUsers(
	rc *middleware.RequestContext,
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	const page uint8 = 1
	tmpl := template.Must(
		template.ParseFiles(
			"template/api-users.html",
			"template/api-users-component.html",
			"template/base.html",
		),
	)

	apiUserRange := db.ApiUserRange{Limit: setting.Setting.PageSize + 1, Offset: 0}
	errs := db.PerformBatch(r.Context(), rc.DbPool, []db.BatchSet{apiUserRange.Get})
	apiUserErr := errs[0]
	if apiUserErr != nil {
		rc.Logger.Error(apiUserErr.Error())
		common.ErrorResp(w)
		return
	}
	data := apiUsersData{
		CreatedState: false,
		Page:         page,
		ApiUsers:     apiUserRange.ApiUsers,
		PostXsrf:     postApiUserXsrf(rc.AdminId),
	}
	if data.nextPageExists() {
		data.ApiUsers = data.ApiUsers[:len(data.ApiUsers)-1]
		data.NextPage = apiUsersApiUrl(page + 1)
	} else {
		data.populateFullPage()
	}
	tmpl.Execute(w, data)
}

func apiUsersCreateApi(
	rc *middleware.RequestContext,
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	data := apiUsersData{
		PostXsrf: postApiUserXsrf(rc.AdminId),
	}
	tmpl := template.Must(template.ParseFiles("template/api-users-component.html")).
		Lookup("api-users-create-form")
	var apiUser db.ApiUser
	err := common.BindJSON(r, &apiUser)
	if err != nil {
		rc.Logger.Error(err.Error())
		common.ErrorResp(w)
		return
	}
	sanitizeApiUser(rc, &apiUser)
	err = rc.Validate.StructPartial(apiUser, "Name")
	if err != nil {
		rc.Logger.Info(err.Error())
		data.InputErr = "Name should be 3-30 characters long"
		tmpl.Execute(w, data)
		return
	}

	rawKey, err := crypto.GenerateUrlSafeString()
	if err != nil {
		rc.Logger.Error(err.Error())
		common.ErrorResp(w)
		return
	}
	hashedKey, err := crypto.HashKey([]byte(rawKey))
	if err != nil {
		rc.Logger.Error(err.Error())
		common.ErrorResp(w)
		return
	}
	apiUser.Key = hashedKey

	errs := db.PerformBatch(r.Context(), rc.DbPool, []db.BatchSet{apiUser.Create})
	apiUserErr := errs[0]
	if apiUserErr != nil {
		errStruct := db.ErrorAsStruct(apiUserErr)
		switch errStruct.(type) {
		case db.UniqueViolation:
			rc.Logger.Info(apiUserErr.Error())
			data.InputErr = "Name exists"
			tmpl.Execute(w, data)
			return
		default:
			rc.Logger.Error(apiUserErr.Error())
			common.ErrorResp(w)
			return
		}
	}
	data.CreatedState = true
	data.Name = apiUser.Name
	data.Key = rawKey
	tmpl.Execute(w, data)
}

func apiUsersApi(
	rc *middleware.RequestContext,
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	data := apiUsersData{
		CreatedState: false,
	}
	tmpl := template.Must(template.ParseFiles("template/api-users-component.html")).
		Lookup("api-users-list")
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		rc.Logger.Info(err.Error())
		tmpl.Execute(w, data)
		return
	}
	data.Page = uint8(page)
	offset := (data.Page - 1) * uint8(setting.Setting.PageSize)
	apiUserRange := db.ApiUserRange{Limit: setting.Setting.PageSize + 1, Offset: offset}
	errs := db.PerformBatch(r.Context(), rc.DbPool, []db.BatchSet{apiUserRange.Get})
	apiUserErr := errs[0]
	if apiUserErr != nil {
		rc.Logger.Error(apiUserErr.Error())
		common.ErrorResp(w)
		return
	}
	data.ApiUsers = apiUserRange.ApiUsers

	if data.nextPageExists() {
		data.ApiUsers = data.ApiUsers[:len(data.ApiUsers)-1]
		data.NextPage = apiUsersApiUrl(data.Page + 1)
	} else {
		data.populateFullPage()
	}
	if data.Page > 1 {
		data.PreviousPage = apiUsersApiUrl(data.Page - 1)
	}

	tmpl.Execute(w, data)
}
