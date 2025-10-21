package main

import (
	"art-collection-system/internal/config"
	"art-collection-system/internal/database"
	"art-collection-system/internal/handler"
	"art-collection-system/internal/middleware"
	"art-collection-system/internal/repository"
	"art-collection-system/internal/routes"
	"art-collection-system/internal/service"
	"art-collection-system/internal/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := initLogger(cfg.Log)
	defer logger.Sync()

	logger.Info("Starting art collection system...")

	// Initialize MySQL database
	logger.Info("Connecting to MySQL...")
	db, err := database.InitMySQL(database.MySQLConfig{
		Host:         cfg.Database.MySQL.Host,
		Port:         cfg.Database.MySQL.Port,
		User:         cfg.Database.MySQL.User,
		Password:     cfg.Database.MySQL.Password,
		DBName:       cfg.Database.MySQL.DBName,
		MaxIdleConns: cfg.Database.MySQL.MaxIdleConns,
		MaxOpenConns: cfg.Database.MySQL.MaxOpenConns,
	})
	if err != nil {
		logger.Fatal("Failed to connect to MySQL", zap.Error(err))
	}
	logger.Info("MySQL connected successfully")

	// Initialize Redis
	logger.Info("Connecting to Redis...")
	redisClient, err := database.InitRedis(database.RedisConfig{
		Host:     cfg.Database.Redis.Host,
		Port:     cfg.Database.Redis.Port,
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.DB,
	})
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Redis connected successfully")

	// Initialize JWT
	utils.InitJWT(cfg.JWT.Secret)
	logger.Info("JWT initialized")

	// Initialize email service
	emailService := utils.NewEmailService(&cfg.Email)
	logger.Info("Email service initialized")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	artworkRepo := repository.NewArtworkRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, redisClient, emailService)
	userService := service.NewUserService(userRepo, artworkRepo)
	activityService := service.NewActivityService(activityRepo)
	fileService := service.NewFileService(cfg.Upload.Path)
	artworkService := service.NewArtworkService(artworkRepo, activityService, fileService)
	adminService := service.NewAdminService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	activityHandler := handler.NewActivityHandler(activityService)
	artworkHandler := handler.NewArtworkHandler(artworkService, fileService)
	adminHandler := handler.NewAdminHandler(artworkService, adminService)

	// Initialize middlewares
	authMiddleware := middleware.AuthMiddleware(authService)
	adminMiddleware := middleware.AdminMiddleware()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create Gin engine
	r := gin.New()

	// Use recovery middleware to handle panics
	r.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(
		r,
		authHandler,
		userHandler,
		activityHandler,
		artworkHandler,
		adminHandler,
		authMiddleware,
		adminMiddleware,
		redisClient,
	)

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(cfg.Upload.Path, 0755); err != nil {
		logger.Fatal("Failed to create uploads directory", zap.Error(err))
	}

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Server starting", zap.String("address", addr), zap.String("mode", cfg.Server.Mode))

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

// initLogger initializes Zap logger based on configuration
func initLogger(logConfig config.LogConfig) *zap.Logger {
	// Parse log level
	var level zapcore.Level
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create core
	var core zapcore.Core
	if logConfig.File != "" {
		// Log to file
		file, err := os.OpenFile(logConfig.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level)

		// Also log to console
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)

		// Combine both cores
		core = zapcore.NewTee(fileCore, consoleCore)
	} else {
		// Log to console only
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	}

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}
