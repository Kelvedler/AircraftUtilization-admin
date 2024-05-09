package auth

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/db"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

const (
	ClaimSubject        = "sub"
	ClaimIssuedAt       = "iat"
	ClaimExpirationTime = "exp"
)

var (
	ErrSubject    = errors.New("Invalid JWT subject")
	ErrUnexpected = errors.New("Unexpected validation error")
	ErrAccess     = errors.New("Subject has no access right")
	ErrIssued     = errors.New("Invalid JWT issued at")
	ErrExpiration = errors.New("Token expired")
)

func JwtExpiration() time.Time {
	return time.Now().Add(
		time.Duration(setting.Setting.Jwt.ExpirationDeltaMinutes) * time.Minute,
	).UTC()
}

func ReissueAllowed(issuedAt int64) bool {
	issuedExp := setting.Setting.Jwt.ExpirationDeltaMinutes * 60
	nowUnix := time.Now().Unix()
	expUnix := time.Unix(int64(issuedExp), int64(0)).Unix() + issuedAt
	return nowUnix < expUnix
}

func IssueJwt(admin db.Admin) (string, error) {
	claims := jwt.MapClaims{}
	claims[ClaimSubject] = admin.ID.String()
	claims[ClaimIssuedAt] = time.Now().Unix()
	claims[ClaimExpirationTime] = JwtExpiration().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(setting.Setting.SecretKey))
}

func ReissueJwt(claims jwt.MapClaims) (string, error) {
	claims[ClaimExpirationTime] = JwtExpiration().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(setting.Setting.SecretKey))
}

func ValidateJwt(
	logger *slog.Logger,
	tokenString string,
) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(setting.Setting.SecretKey), nil
	})
	if err != nil {
		logger.Info(err.Error())
		return jwt.MapClaims{}, false
	}
	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error(ErrUnexpected.Error())
		return jwt.MapClaims{}, false
	}
	_, ok = tokenClaims[ClaimSubject].(string)
	if !ok {
		logger.Error(ErrSubject.Error())
		return jwt.MapClaims{}, false
	}
	_, ok = tokenClaims[ClaimIssuedAt].(float64)
	if !ok {
		logger.Error(ErrIssued.Error())
		return jwt.MapClaims{}, false
	}
	tsNow := time.Now().Unix()
	expiredAt, ok := tokenClaims[ClaimExpirationTime].(float64)
	if !ok || tsNow > int64(expiredAt) {
		logger.Info(ErrExpiration.Error())
		return jwt.MapClaims{}, false
	}

	return tokenClaims, true
}
