package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/common"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/db"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/view"
)

func main() {
	ctx := context.Background()
	setting.InitSetting()
	mainLogger := common.MainLogger()
	validate := common.NewValidator()
	sanitize := common.NewSanitizer()
	dbpool := db.NewConnectionPool(ctx, mainLogger)
	router := view.BaseRouter(dbpool, sanitize, validate, mainLogger)
	http.ListenAndServe(fmt.Sprintf(":%d", setting.Setting.Port), router)
}
