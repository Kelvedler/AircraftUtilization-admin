package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/db"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

const cookiePath = "/"

func setTokenCookie(w http.ResponseWriter, token string) error {
	jwtEnv := setting.Setting.Jwt
	cookie := http.Cookie{
		Name:     "access",
		Value:    token,
		Path:     cookiePath,
		Domain:   jwtEnv.Domain,
		Secure:   jwtEnv.SecureCookies,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  JwtExpiration(),
	}
	http.SetCookie(w, &cookie)
	return nil
}

func SetNewTokenCookie(w http.ResponseWriter, admin db.Admin) error {
	token, err := IssueJwt(admin)
	if err != nil {
		return err
	}
	return setTokenCookie(w, token)
}

func ReissueTokenCookie(w http.ResponseWriter, claims jwt.MapClaims) error {
	token, err := ReissueJwt(claims)
	if err != nil {
		return err
	}
	return setTokenCookie(w, token)
}
