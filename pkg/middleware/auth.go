package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/auth"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/db"
)

func PerformAuth(
	logger *slog.Logger,
	dbpool *pgxpool.Pool,
	w http.ResponseWriter,
	r *http.Request,
) (adminId uuid.UUID, returnErr error) {
	returnErr = errors.New("Unauthorized")

	accessToken, err := r.Cookie("access")
	if err != nil {
		return adminId, returnErr
	}

	tokenClaims, ok := auth.ValidateJwt(logger, accessToken.Value)
	if !ok {
		return adminId, returnErr
	}

	adminIdStr := tokenClaims[auth.ClaimSubject].(string)
	adminId, err = uuid.Parse(adminIdStr)
	if err != nil {
		return uuid.Nil, returnErr
	}

	if auth.ReissueAllowed(int64(tokenClaims[auth.ClaimIssuedAt].(float64))) {
		err = auth.ReissueTokenCookie(w, tokenClaims)
		if err != nil {
			return adminId, err
		}
		return adminId, nil
	}

	caller := db.Admin{ID: adminId}
	errs := db.PerformBatch(r.Context(), dbpool, []db.BatchSet{caller.GetById})
	adminErr := errs[0]
	if adminErr == nil {
		auth.SetNewTokenCookie(w, caller)
		return adminId, nil
	}

	errStruct := db.ErrorAsStruct(adminErr)
	switch errStruct.(type) {
	case db.DoesNotExist:
		logger.Info("Not found")
		return uuid.Nil, returnErr
	case db.ContextCanceled:
		logger.Warn("Context canceled")
		return uuid.Nil, returnErr
	default:
		panic(fmt.Sprintf("Unexpected err type, %t", errStruct))
	}
}

func UnauthorizedRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
