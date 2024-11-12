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
	roleRepository       repository.RoleRepository       = repository.NewRoleRepository(db)
	menuRepository       repository.MenuRepository       = repository.NewMenuRepository(db)

	jwtService        service.JWTService        = service.NewJwtService()
	authService       service.AuthService       = service.NewAuthService(userRepository)
	userService       service.UserService       = service.NewUserService(userRepository)
	perusahaanService service.PerusahaanService = service.NewPerusahaanService(perusahaanRepository, validate)
	roleService       service.RoleService       = service.NewRoleService(roleRepository, validate)
	menuService       service.MenuService       = service.NewMenuService(menuRepository, validate)

	authController       controller.AuthController       = controller.NewAuthController(authService, jwtService)
	userController       controller.UserController       = controller.NewUserController(userService, jwtService)
	perusahaanController controller.PerusahaanController = controller.NewPerusahaanController(perusahaanService, jwtService)
	roleController       controller.RoleController       = controller.NewRoleController(roleService, jwtService)
	menuController       controller.MenuController       = controller.NewMenuController(menuService, jwtService)
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

	roleRoutes := r.Group("api/role")
	{
		roleRoutes.POST("/", roleController.CreateRole)
		roleRoutes.PATCH("/:roleId", roleController.UpdateRole)
		roleRoutes.GET("/", roleController.FindRoleAll)
		roleRoutes.GET("/:roleId", roleController.FindRoleById)
		roleRoutes.DELETE("/:roleId", roleController.DeleteRole)
	}

	menuRoutes := r.Group("api/role")
	{
		menuRoutes.POST("/", menuController.CreateMenu)
		menuRoutes.PATCH("/:roleId", menuController.UpdateMenu)
		menuRoutes.GET("/", menuController.FindMenuAll)
		menuRoutes.GET("/:roleId", menuController.FindMenuById)
		menuRoutes.DELETE("/:roleId", menuController.DeleteMenu)
	}

	return r
}
