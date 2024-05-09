package setting

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type AppSetting struct {
	PageSize      uint8
	Argon2        Argon2
	ApiKey        ApiKey
	SecretKey     string
	LogLevel      slog.Level
	AdminPassword string
	DatabaseUrl   string
	Jwt           Jwt
}

type Argon2 struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint8
	KeyLength   uint32
}

type ApiKey struct {
	Length uint8
}

type Jwt struct {
	Domain                 string
	SecureCookies          bool
	ExpirationDeltaMinutes int
}

var Setting AppSetting

func setPageSize() {
	Setting.PageSize = 20
}

func setArgon2() {
	Setting.Argon2 = Argon2{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func setApiKey() {
	Setting.ApiKey = ApiKey{Length: 32}
}

func ensureValueExists(key, value string, logger *slog.Logger) {
	if value == "" {
		logger.Error(fmt.Sprintf("Could not get '%s'", key))
		os.Exit(1)
	}
}

func setSecretKey(logger *slog.Logger) {
	envKey := "SECRET_KEY"
	secretKey := os.Getenv(envKey)
	ensureValueExists(envKey, secretKey, logger)
	Setting.SecretKey = secretKey
}

func setLogLevel(logger *slog.Logger) {
	envKey := "LOG_LEVEL"
	level := os.Getenv(envKey)
	switch level {
	case "DEBUG":
		Setting.LogLevel = slog.LevelDebug
	case "WARN":
		Setting.LogLevel = slog.LevelWarn
	case "ERROR":
		Setting.LogLevel = slog.LevelError
	default:
		if level != "INFO" {
			logger.Info(fmt.Sprintf("Could not get '%s', set to INFO", envKey))
		}
		Setting.LogLevel = slog.LevelInfo
	}
}

func setDatabaseUrl(logger *slog.Logger) {
	envDbUrlKey := "DATABASE_URL"
	databaseUrl := os.Getenv(envDbUrlKey)
	ensureValueExists(envDbUrlKey, databaseUrl, logger)

	Setting.DatabaseUrl = databaseUrl
}

func setJwt(logger *slog.Logger) {
	secure, err := strconv.ParseBool(os.Getenv("JWT_SECURE_COOKIES"))
	if err != nil {
		logger.Error("Could not get 'JWT_SECURE_COOKIES'")
		os.Exit(1)
	}
	jwtDomainKey := "JWT_DOMAIN"
	domain := os.Getenv(jwtDomainKey)
	ensureValueExists(jwtDomainKey, domain, logger)

	exp, err := strconv.Atoi(os.Getenv("JWT_EXP_DELTA_MINUTES"))
	if err != nil {
		logger.Error("Could not get 'JWT_EXP_DELTA_MINUTES'")
		os.Exit(1)
	}
	jwt := Jwt{
		Domain:                 domain,
		SecureCookies:          secure,
		ExpirationDeltaMinutes: exp,
	}
	Setting.Jwt = jwt
}

func InitSetting() {
	setPageSize()
	setArgon2()
	setApiKey()

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
	).With(slog.String("process", "main"))

	setSecretKey(logger)
	setLogLevel(logger)
	setDatabaseUrl(logger)
	setJwt(logger)
}
