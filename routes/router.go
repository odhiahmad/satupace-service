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
	validate                *validator.Validate                = validator.New()
	db                      *gorm.DB                           = config.SetupDatabaseConnection()
	userRepository          repository.UserRepository          = repository.NewUserRepository(db)
	roleRepository          repository.RoleRepository          = repository.NewRoleRepository(db)
	businessTypeRepository  repository.BusinessTypeRepository  = repository.NewBusinessTypeRepository(db)
	paymentMethodRepository repository.PaymentMethodRepository = repository.NewPaymentMethodRepository(db)
	productUnitRepository   repository.ProductUnitRepository   = repository.NewProductUnitRepository(db)

	jwtService           service.JWTService           = service.NewJwtService()
	authService          service.AuthService          = service.NewAuthService(userRepository)
	userService          service.UserService          = service.NewUserService(userRepository, validate)
	roleService          service.RoleService          = service.NewRoleService(roleRepository, validate)
	businessTypeService  service.BusinessTypeService  = service.NewBusinessTypeService(businessTypeRepository, validate)
	paymentMethodService service.PaymentMethodService = service.NewPaymentMethodService(paymentMethodRepository, validate)
	productUnitService   service.ProductUnitService   = service.NewProductUnitService(productUnitRepository, validate)

	authController          controller.AuthController          = controller.NewAuthController(authService)
	userController          controller.UserController          = controller.NewUserController(userService, jwtService)
	roleController          controller.RoleController          = controller.NewRoleController(roleService, jwtService)
	businessTypeController  controller.BusinessTypeController  = controller.NewBusinessTypeController(businessTypeService, jwtService)
	paymentMethodController controller.PaymentMethodController = controller.NewPaymentMethodController(paymentMethodService, jwtService)
	productUnitController   controller.ProductUnitController   = controller.NewProductUnitController(productUnitService, jwtService)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
	}

	userRoutes := r.Group("api/user")
	{
		userRoutes.POST("", userController.InsertRegistration)
		userRoutes.POST("/create", userController.CreateUser)
		userRoutes.PUT("/update", userController.UpdateUser)
	}

	roleRoutes := r.Group("api/role")
	{
		roleRoutes.POST("/", roleController.CreateRole)
		roleRoutes.PATCH("/:roleId", roleController.UpdateRole)
		roleRoutes.GET("/", roleController.FindRoleAll)
		roleRoutes.GET("/:roleId", roleController.FindRoleById)
		roleRoutes.DELETE("/:roleId", roleController.DeleteRole)
	}

	businessTypeRoutes := r.Group("api/business-type")
	{
		businessTypeRoutes.POST("/", businessTypeController.CreateBusinessType)
		businessTypeRoutes.PATCH("/:businessTypeId", businessTypeController.UpdateBusinessType)
		businessTypeRoutes.GET("/", businessTypeController.FindBusinessTypeAll)
		businessTypeRoutes.GET("/:businessTypeId", businessTypeController.FindBusinessTypeById)
		businessTypeRoutes.DELETE("/:businessTypeId", businessTypeController.DeleteBusinessType)
	}

	paymentMethodRoutes := r.Group("api/business-type")
	{
		paymentMethodRoutes.POST("/", paymentMethodController.CreatePaymentMethod)
		paymentMethodRoutes.PATCH("/:paymentMethodId", paymentMethodController.UpdatePaymentMethod)
		paymentMethodRoutes.GET("/", paymentMethodController.FindPaymentMethodAll)
		paymentMethodRoutes.GET("/:paymentMethodId", paymentMethodController.FindPaymentMethodById)
		paymentMethodRoutes.DELETE("/:paymentMethodId", paymentMethodController.DeletePaymentMethod)
	}

	productUnitRoutes := r.Group("api/business-type")
	{
		productUnitRoutes.POST("/", productUnitController.CreateProductUnit)
		productUnitRoutes.PATCH("/:productUnitId", productUnitController.UpdateProductUnit)
		productUnitRoutes.GET("/", productUnitController.FindProductUnitAll)
		productUnitRoutes.GET("/:productUnitId", productUnitController.FindProductUnitById)
		productUnitRoutes.DELETE("/:productUnitId", productUnitController.DeleteProductUnit)
	}

	return r
}
