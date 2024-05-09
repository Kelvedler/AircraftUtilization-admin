package crypto

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

func encodeHash(salt, hash []byte, p setting.Argon2) string {
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.Memory,
		p.Iterations,
		p.Parallelism,
		b64Salt,
		b64Hash,
	)
}

func decodeHash(encodedHash string) (p *setting.Argon2, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}
	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}
	p = &setting.Argon2{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}
	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint8(len(salt))
	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))
	return p, salt, hash, nil
}

func HashKey(password []byte) (encodedHash string, err error) {
	p := setting.Setting.Argon2
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey(
		password,
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)
	return encodeHash(salt, hash, p), nil
}

func CompareKeys(rawKey, hashedKey string) (identical bool, err error) {
	p, salt, hash, err := decodeHash(hashedKey)
	if err != nil {
		return false, err
	}
	otherHash := argon2.IDKey(
		[]byte(rawKey),
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
