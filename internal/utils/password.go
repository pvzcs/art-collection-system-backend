package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 Bcrypt 加密密码，cost 为 12
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword 比较明文密码和哈希密码是否匹配
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
