package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

const (
	// MaxFileSize 最大文件大小 10MB
	MaxFileSize = 10 * 1024 * 1024 // 10MB in bytes
)

var (
	// AllowedImageExtensions 允许的图片扩展名
	AllowedImageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}

	// ImageMagicNumbers 图片文件的魔数（文件头标识）
	ImageMagicNumbers = map[string][]byte{
		"jpeg": {0xFF, 0xD8, 0xFF},
		"png":  {0x89, 0x50, 0x4E, 0x47},
		"gif":  {0x47, 0x49, 0x46, 0x38},
		"webp": {0x52, 0x49, 0x46, 0x46}, // RIFF (WebP 的前4字节)
		"bmp":  {0x42, 0x4D},
	}

	ErrFileTooLarge      = errors.New("文件大小超过限制（最大 10MB）")
	ErrInvalidFileType   = errors.New("不支持的文件类型，仅允许图片格式")
	ErrInvalidFileHeader = errors.New("文件内容与扩展名不匹配")
)

// ValidateImageFile 验证上传的图片文件
// 检查文件大小、扩展名和文件内容
// Requirements: 11.1
func ValidateImageFile(fileHeader *multipart.FileHeader) error {
	// 1. 检查文件大小
	if fileHeader.Size > MaxFileSize {
		return ErrFileTooLarge
	}

	// 2. 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isAllowedExtension(ext) {
		return ErrInvalidFileType
	}

	// 3. 验证文件内容（防止伪造扩展名）
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 读取文件头（前 512 字节足够识别大多数文件类型）
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("无法读取文件内容: %w", err)
	}
	buffer = buffer[:n]

	// 验证文件魔数
	if !isValidImageContent(buffer) {
		return ErrInvalidFileHeader
	}

	return nil
}

// isAllowedExtension 检查文件扩展名是否允许
func isAllowedExtension(ext string) bool {
	for _, allowed := range AllowedImageExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

// isValidImageContent 通过文件魔数验证文件内容
func isValidImageContent(data []byte) bool {
	if len(data) < 2 {
		return false
	}

	// 检查各种图片格式的魔数
	for _, magic := range ImageMagicNumbers {
		if bytes.HasPrefix(data, magic) {
			return true
		}
	}

	// WebP 需要额外检查（RIFF 后面应该跟着 WEBP）
	if len(data) >= 12 && bytes.HasPrefix(data, []byte{0x52, 0x49, 0x46, 0x46}) {
		if bytes.Equal(data[8:12], []byte{0x57, 0x45, 0x42, 0x50}) {
			return true
		}
	}

	return false
}

// FormatFileSize 格式化文件大小为人类可读的格式
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
