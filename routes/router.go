package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/config"
	"github.com/odhiahmad/kasirku-service/controller"
	"github.com/odhiahmad/kasirku-service/repository"
	"github.com/odhiahmad/kasirku-service/service"
	"gorm.io/gorm"
)

var (
	// Initialize the validator
	validate             *validator.Validate             = validator.New()
	db                   *gorm.DB                        = config.SetupDatabaseConnection()
	userRepository       repository.UserRepository       = repository.NewUserRepository(db)
	perusahaanRepository repository.PerusahaanRepository = repository.NewPerusahaanRepository(db)

	jwtService        service.JWTService        = service.NewJwtService()
	authService       service.AuthService       = service.NewAuthService(userRepository)
	userService       service.UserService       = service.NewUserService(userRepository)
	perusahaanService service.PerusahaanService = service.NewPerusahaanService(perusahaanRepository, validate)

	authController       controller.AuthController       = controller.NewAuthController(authService, jwtService)
	userController       controller.UserController       = controller.NewUserController(userService, jwtService)
	perusahaanController controller.PerusahaanController = controller.NewPerusahaanController(perusahaanService, jwtService)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
	}

	userRoutes := r.Group("api/user")
	{
		userRoutes.POST("/create", userController.CreateUser)
		userRoutes.PUT("/update", userController.UpdateUser)
	}

	perusahaanRoutes := r.Group("api/perusahaan")
	{
		perusahaanRoutes.POST("/", perusahaanController.CreatePerusahaan)
		perusahaanRoutes.PATCH("/:perusahaanId", perusahaanController.UpdatePerusahaan)
		perusahaanRoutes.GET("/", perusahaanController.FindPerusahaanAll)
		perusahaanRoutes.GET("/:perusahaanId", perusahaanController.FindPerusahaanById)
		perusahaanRoutes.DELETE("/:perusahaanId", perusahaanController.DeletePerusahaan)
	}
	return r
}
