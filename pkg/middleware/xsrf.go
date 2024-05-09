package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/net/xsrftoken"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

func ValidateForXSRF(r *http.Request, adminId uuid.UUID) bool {
	if r.Method == http.MethodGet {
		return true
	}
	token := r.Header.Get("_xsrf")
	return xsrftoken.Valid(token, setting.Setting.SecretKey, adminId.String(), r.URL.String())
}
