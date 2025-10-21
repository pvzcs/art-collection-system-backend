package routes

import (
	"art-collection-system/internal/handler"
	"art-collection-system/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	activityHandler *handler.ActivityHandler,
	artworkHandler *handler.ArtworkHandler,
	adminHandler *handler.AdminHandler,
	authMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc,
	redisClient *redis.Client,
) {
	// Apply global middlewares
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// API v1 routes
	v1 := r.Group("/api/v1")

	// Public routes (no authentication required)
	setupPublicRoutes(v1, authHandler, activityHandler, redisClient)

	// Protected routes (authentication required)
	setupProtectedRoutes(v1, authHandler, userHandler, activityHandler, artworkHandler, authMiddleware, redisClient)

	// Admin routes (authentication + admin role required)
	setupAdminRoutes(v1, activityHandler, adminHandler, authMiddleware, adminMiddleware)
}

// setupPublicRoutes configures public routes
func setupPublicRoutes(
	rg *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	activityHandler *handler.ActivityHandler,
	redisClient *redis.Client,
) {
	// Authentication routes
	auth := rg.Group("/auth")
	{
		auth.POST("/send-code", middleware.VerificationCodeRateLimiter(redisClient), authHandler.SendVerificationCode)
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", middleware.LoginRateLimiter(redisClient), authHandler.Login)
	}

	// Public activity routes
	activities := rg.Group("/activities")
	{
		activities.GET("", activityHandler.ListActivities)
		activities.GET("/:id", activityHandler.GetActivity)
	}
}

// setupProtectedRoutes configures routes that require authentication
func setupProtectedRoutes(
	rg *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	activityHandler *handler.ActivityHandler,
	artworkHandler *handler.ArtworkHandler,
	authMiddleware gin.HandlerFunc,
	redisClient *redis.Client,
) {
	protected := rg.Group("")
	protected.Use(authMiddleware)

	// Authentication routes (logout requires auth)
	auth := protected.Group("/auth")
	{
		auth.POST("/logout", authHandler.Logout)
	}

	// User routes
	user := protected.Group("/user")
	{
		user.GET("/profile", userHandler.GetProfile)
		user.PUT("/profile", userHandler.UpdateProfile)
		user.PUT("/password", userHandler.ChangePassword)
	}

	// User artworks (personal space)
	users := protected.Group("/users")
	{
		users.GET("/:id/artworks", userHandler.GetUserArtworks)
	}

	// Artwork routes
	artworks := protected.Group("/artworks")
	{
		artworks.POST("", middleware.UploadRateLimiter(redisClient), artworkHandler.UploadArtwork)
		artworks.DELETE("/:id", artworkHandler.DeleteArtwork)
		artworks.GET("/:id", artworkHandler.GetArtwork)
		artworks.GET("/:id/image", artworkHandler.ServeImage)
	}
}

// setupAdminRoutes configures routes that require admin role
func setupAdminRoutes(
	rg *gin.RouterGroup,
	activityHandler *handler.ActivityHandler,
	adminHandler *handler.AdminHandler,
	authMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc,
) {
	admin := rg.Group("/admin")
	admin.Use(authMiddleware, adminMiddleware)

	// Activity management
	activities := admin.Group("/activities")
	{
		activities.POST("", activityHandler.CreateActivity)
		activities.PUT("/:id", activityHandler.UpdateActivity)
		activities.DELETE("/:id", activityHandler.DeleteActivity)
	}

	// Artwork review
	admin.GET("/review-queue", adminHandler.GetReviewQueue)
	artworks := admin.Group("/artworks")
	{
		artworks.PUT("/:id/review", adminHandler.ReviewArtwork)
		artworks.PUT("/batch-review", adminHandler.BatchReviewArtworks)
	}

	// User management
	users := admin.Group("/users")
	{
		users.GET("", adminHandler.ListUsers)
		users.PUT("/:id/role", adminHandler.UpdateUserRole)
		users.GET("/:id/statistics", adminHandler.GetUserStatistics)
	}
}
