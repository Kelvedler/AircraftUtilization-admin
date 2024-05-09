package view

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/microcosm-cc/bluemonday"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/middleware"
)

func staticFilepath(mainLogger *slog.Logger) string {
	wd, err := os.Getwd()
	if err != nil {
		mainLogger.Error(err.Error())
		os.Exit(1)
	}
	return wd + "/static/"
}

func BaseRouter(
	dbpool *pgxpool.Pool,
	sanitize *bluemonday.Policy,
	validate *validator.Validate,
	mainLogger *slog.Logger,
) *httprouter.Router {
	router := httprouter.New()
	handlerContext := middleware.NewHandlerContext(dbpool, sanitize, validate)

	router.ServeFiles("/static/*filepath", http.Dir(staticFilepath(mainLogger)))

	router.GET("/", middleware.Unrestricted.Wrapper(index, handlerContext))
	router.GET("/api-users", middleware.AuthOnly.Wrapper(apiUsers, handlerContext))

	router.POST("/api/v1/sign-in", middleware.Unrestricted.Wrapper(signInApi, handlerContext))
	router.GET("/api/v1/api-users/", middleware.AuthOnlyApi.Wrapper(apiUsersApi, handlerContext))
	router.POST(
		"/api/v1/api-users",
		middleware.AuthOnlyApi.Wrapper(apiUsersCreateApi, handlerContext),
	)
	return router
}
