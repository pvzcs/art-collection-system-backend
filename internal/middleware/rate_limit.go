package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	MaxRequests int           // 最大请求数
	Window      time.Duration // 时间窗口
	KeyPrefix   string        // Redis key 前缀
}

// RateLimiter 速率限制器
type RateLimiter struct {
	redis  *redis.Client
	config RateLimitConfig
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(redis *redis.Client, config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		redis:  redis,
		config: config,
	}
}

// Middleware 速率限制中间件
func (rl *RateLimiter) Middleware(keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := rl.config.KeyPrefix + keyFunc(c)
		ctx := context.Background()

		// 获取当前计数
		count, err := rl.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(500, gin.H{"code": 500, "message": "服务器内部错误"})
			c.Abort()
			return
		}

		// 检查是否超过限制
		if count >= rl.config.MaxRequests {
			c.JSON(429, gin.H{"code": 429, "message": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}

		// 增加计数
		pipe := rl.redis.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, rl.config.Window)
		_, err = pipe.Exec(ctx)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "message": "服务器内部错误"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// VerificationCodeRateLimiter 验证码发送速率限制（每个邮箱每分钟 1 次）
func VerificationCodeRateLimiter(redis *redis.Client) gin.HandlerFunc {
	limiter := NewRateLimiter(redis, RateLimitConfig{
		MaxRequests: 1,
		Window:      time.Minute,
		KeyPrefix:   "rate_limit:verify_code:",
	})

	return limiter.Middleware(func(c *gin.Context) string {
		// 从查询参数或表单中获取邮箱
		email := c.Query("email")
		if email == "" {
			email = c.PostForm("email")
		}
		// 如果还是空，尝试从 JSON body 中获取
		if email == "" {
			var req struct {
				Email string `json:"email"`
			}
			// 读取 body 内容
			bodyBytes, _ := c.GetRawData()
			if len(bodyBytes) > 0 {
				// 重新设置 body 以便后续处理器可以读取
				c.Request.Body = &readCloser{bytes.NewReader(bodyBytes)}
				// 尝试解析 JSON
				json.Unmarshal(bodyBytes, &req)
				email = req.Email
			}
		}
		return email
	})
}

// readCloser 实现 io.ReadCloser 接口
type readCloser struct {
	*bytes.Reader
}

func (rc *readCloser) Close() error {
	return nil
}

// LoginRateLimiter 登录尝试速率限制（每个 IP 每分钟 5 次）
func LoginRateLimiter(redis *redis.Client) gin.HandlerFunc {
	limiter := NewRateLimiter(redis, RateLimitConfig{
		MaxRequests: 5,
		Window:      time.Minute,
		KeyPrefix:   "rate_limit:login:",
	})

	return limiter.Middleware(func(c *gin.Context) string {
		return c.ClientIP()
	})
}

// UploadRateLimiter 文件上传速率限制（每个用户每分钟 10 次）
func UploadRateLimiter(redis *redis.Client) gin.HandlerFunc {
	limiter := NewRateLimiter(redis, RateLimitConfig{
		MaxRequests: 10,
		Window:      time.Minute,
		KeyPrefix:   "rate_limit:upload:",
	})

	return limiter.Middleware(func(c *gin.Context) string {
		userID, exists := c.Get("user_id")
		if !exists {
			return ""
		}
		return fmt.Sprintf("%v", userID)
	})
}
