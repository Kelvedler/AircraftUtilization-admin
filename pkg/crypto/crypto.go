package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

func generateRandomBytes(n uint8) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateUrlSafeSecret() (string, error) {
	b, err := generateRandomBytes(setting.Setting.ApiKey.SecretLength)
	if err != nil {
		return "", err
	}
	encodedStr := base64.RawURLEncoding.EncodeToString(b)
	return encodedStr, nil
}
