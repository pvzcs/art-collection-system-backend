package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// 此工具用于生成 bcrypt 密码哈希
// 使用方法: go run scripts/generate_password.go <password>
// 示例: go run scripts/generate_password.go Admin123456

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run scripts/generate_password.go <password>")
		fmt.Println("示例: go run scripts/generate_password.go Admin123456")
		os.Exit(1)
	}

	password := os.Args[1]

	// 使用 bcrypt cost=12 生成密码哈希
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		fmt.Printf("生成密码哈希失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("密码哈希生成成功！")
	fmt.Println("原始密码:", password)
	fmt.Println("密码哈希:", string(hashedBytes))
	fmt.Println("\n您可以在 SQL 脚本中使用此哈希值。")
}
