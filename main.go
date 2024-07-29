package main

import (
	"fmt"

	"lazy-auth/app/handler"
	"lazy-auth/app/middleware"
	repository "lazy-auth/app/repository"
	service "lazy-auth/app/service"
	"lazy-auth/common"
	"lazy-auth/config"
	"lazy-auth/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	config := config.ConfigService()
	db := database.InitDatabase(config)

	roleRepository := repository.NewRoleRepository(db)
	sessionRepository := repository.NewSessionRepository(db)
	userRepository := repository.NewUserRepository(db)

	authService := service.NewAuthService(
		userRepository,
		roleRepository,
		sessionRepository,
		config,
	)
	userService := service.NewUserService(userRepository, roleRepository)

	secretGuard := middleware.NewSecretGuard(config)
	tokenGuard := middleware.NewTokenGuard(sessionRepository, config)
	roleGuard := middleware.NewRoleGuard(userRepository)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	if config.Stage == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", common.PasswordValidate)
	}

	api := r.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)

		// Role
		api.GET("/roles", authHandler.GetRoles)
		api.POST("/roles", authHandler.CreateRole)

		// Auth
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.RefreshToken)
		api.POST("/auth/logout", authHandler.Logout)
		api.POST("/auth/forgot-password", authHandler.ForgotPassword)
		api.POST("/auth/reset-password", authHandler.ResetPassword)

		// User
		api.GET(
			"/users",
			tokenGuard.ValidateToken(),
			roleGuard.ValidateRole("admin"),
			userHandler.GetUsers,
		)
		api.POST("/users", userHandler.CreateUser)
		api.POST("/users/admin", secretGuard.ValidateSecret(), userHandler.CreateUserAdmin)
		api.GET("/users/me", tokenGuard.ValidateToken(), userHandler.GetMe)
		api.PATCH("/users/me", tokenGuard.ValidateToken(), userHandler.UpdateMe)
		api.POST("/users/change-password", tokenGuard.ValidateToken(), userHandler.ChangePassword)
	}

	r.Run(fmt.Sprintf(":%v", config.Port))
}
