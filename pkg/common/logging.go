package common

import (
	"log/slog"
	"os"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

func MainLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: setting.Setting.LogLevel}),
	).With(slog.String("process", "main"))
}
