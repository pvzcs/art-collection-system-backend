package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// 此工具用于验证密码哈希是否正确
// 使用方法: go run scripts/verify_password.go <password> <hash>
// 示例: go run scripts/verify_password.go Admin123456 '$2a$12$...'

func main() {
	if len(os.Args) < 3 {
		fmt.Println("使用方法: go run scripts/verify_password.go <password> <hash>")
		fmt.Println("示例: go run scripts/verify_password.go Admin123456 '$2a$12$...'")
		os.Exit(1)
	}

	password := os.Args[1]
	hash := os.Args[2]

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("❌ 密码验证失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 密码验证成功！")
	fmt.Println("密码:", password)
	fmt.Println("哈希:", hash)
}
