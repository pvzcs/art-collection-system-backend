package handler

import (
	"art-collection-system/internal/service"
	"art-collection-system/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new authentication handler instance
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SendVerificationCodeRequest represents the request body for sending verification code
type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

// SendVerificationCode handles sending verification code to email
// POST /api/v1/auth/send-code
func (h *AuthHandler) SendVerificationCode(c *gin.Context) {
	var req SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Validate email format
	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Send verification code
	if err := h.authService.SendVerificationCode(req.Email); err != nil {
		utils.Error(c, 500, "发送验证码失败")
		return
	}

	utils.Success(c, gin.H{"message": "验证码已发送"})
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Validate email format
	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Validate password strength
	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Validate nickname
	if err := utils.ValidateNickname(req.Nickname); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Validate verification code (6 digits)
	if len(req.Code) != 6 {
		utils.Error(c, 400, "验证码格式不正确")
		return
	}

	// Register user
	user, err := h.authService.Register(req.Email, req.Code, req.Password, req.Nickname)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			utils.Error(c, 400, "邮箱已被注册")
		} else if strings.Contains(err.Error(), "verification code") {
			utils.Error(c, 400, err.Error())
		} else {
			utils.Error(c, 500, "注册失败")
		}
		return
	}

	// Return user info (without password)
	utils.Success(c, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"nickname": user.Nickname,
		"role":     user.Role,
	})
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Validate email format
	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Login
	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid email or password") {
			utils.Error(c, 401, "邮箱或密码错误")
		} else {
			utils.Error(c, 500, "登录失败")
		}
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"nickname": user.Nickname,
			"role":     user.Role,
		},
	})
}

// Logout handles user logout
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.Error(c, 401, "未提供认证令牌")
		return
	}

	// Remove "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		utils.Error(c, 401, "令牌格式不正确")
		return
	}

	// Logout (add token to blacklist)
	if err := h.authService.Logout(token); err != nil {
		utils.Error(c, 500, "登出失败")
		return
	}

	utils.Success(c, gin.H{"message": "登出成功"})
}
