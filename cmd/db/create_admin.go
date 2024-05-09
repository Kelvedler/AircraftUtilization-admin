package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/common"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/crypto"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/db"
	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

func main() {
	password := common.PasswordFromCli()
	hashedPassword, err := crypto.HashKey(password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	admin := db.Admin{
		Name:     "Admin",
		Password: hashedPassword,
	}
	ctx := context.Background()
	setting.InitSetting()
	mainLogger := common.MainLogger()
	dbpool := db.NewConnectionPool(ctx, mainLogger)
	errs := db.PerformBatch(ctx, dbpool, []db.BatchSet{admin.Create})
	adminErr := errs[0]
	if adminErr != nil {
		fmt.Println(adminErr)
	} else {
		fmt.Println("Admin created successfuly")
	}
	os.Exit(0)
}
