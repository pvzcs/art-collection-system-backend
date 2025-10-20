package utils

import (
	"errors"
	"regexp"
	"unicode"
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}

	// 邮箱正则表达式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式不正确")
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}

	if len(password) < 8 {
		return errors.New("密码长度至少为 8 位")
	}

	// 检查是否包含至少一个字母和一个数字
	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		return errors.New("密码必须包含字母和数字")
	}

	return nil
}

// ValidateNickname 验证昵称
func ValidateNickname(nickname string) error {
	if nickname == "" {
		return errors.New("昵称不能为空")
	}

	if len(nickname) < 2 || len(nickname) > 50 {
		return errors.New("昵称长度必须在 2-50 个字符之间")
	}

	return nil
}
