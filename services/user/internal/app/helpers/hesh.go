package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// Hesh ...
// Протокол хеша
type Hesh struct{}

// HashPassword ...
// захешировать пароль
func (h *Hesh) HashPassword(password string, salt string) string {

	hesh := hmac.New(sha256.New, []byte(salt))
	hesh.Write([]byte(password))

	return base64.StdEncoding.EncodeToString(hesh.Sum(nil))
}

// CheckPasswordHash ...
// Сравнить пароль
// Сравнивает хеш пароля в базе с хешем
// Поступившего пароля
func (h *Hesh) CheckPasswordHash(password, hash string, salt string) bool {
	expected := h.HashPassword(password, salt)
	return hmac.Equal([]byte(hash), []byte(expected))
}
