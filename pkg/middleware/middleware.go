package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/microcosm-cc/bluemonday"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/common"
)

type Settings struct {
	AuthRequired bool
	XsrfExempt   bool
}

var Unrestricted = Settings{
	AuthRequired: false,
	XsrfExempt:   true,
}

var AuthOnly = Settings{
	AuthRequired: true,
	XsrfExempt:   true,
}

var AuthOnlyApi = Settings{
	AuthRequired: true,
	XsrfExempt:   false,
}

type HandlerContext struct {
	dbpool   *pgxpool.Pool
	sanitize *bluemonday.Policy
	validate *validator.Validate
}

func NewHandlerContext(
	dbpool *pgxpool.Pool,
	sanitize *bluemonday.Policy,
	validate *validator.Validate,
) *HandlerContext {
	return &HandlerContext{
		dbpool:   dbpool,
		sanitize: sanitize,
		validate: validate,
	}
}

type RequestContext struct {
	Logger   *slog.Logger
	AdminId  uuid.UUID
	DbPool   *pgxpool.Pool
	Sanitize *bluemonday.Policy
	Validate *validator.Validate
}

type Handle func(*RequestContext, http.ResponseWriter, *http.Request, httprouter.Params)

func (settings Settings) Wrapper(
	handler Handle,
	handlerContext *HandlerContext,
) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		logger := NewRequestLogger()
		RequestLogger(logger, r)
		var adminId uuid.UUID
		var err error
		if settings.AuthRequired == true {
			adminId, err = PerformAuth(logger, handlerContext.dbpool, w, r)
			if err != nil {
				logger.Info(err.Error())
				UnauthorizedRedirect(w, r)
				return
			}
		}
		if !settings.XsrfExempt && !ValidateForXSRF(r, adminId) {
			common.ErrorResp(w)
			logger.Warn("XSRF token invalid")
			return
		}
		rc := &RequestContext{
			Logger:   logger,
			AdminId:  adminId,
			DbPool:   handlerContext.dbpool,
			Sanitize: handlerContext.sanitize,
			Validate: handlerContext.validate,
		}
		handler(rc, w, r, p)
	}
}
