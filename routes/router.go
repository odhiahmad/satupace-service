package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/config"
	"github.com/odhiahmad/kasirku-service/controller"
	"github.com/odhiahmad/kasirku-service/middleware"
	"github.com/odhiahmad/kasirku-service/repository"
	"github.com/odhiahmad/kasirku-service/service"
	"gorm.io/gorm"
)

var (
	// Initialize the validator
	validate                  *validator.Validate                  = validator.New()
	db                        *gorm.DB                             = config.SetupDatabaseConnection()
	userRepository            repository.UserRepository            = repository.NewUserRepository(db)
	userBusinessRepository    repository.UserBusinessRepository    = repository.NewUserBusinessRepository(db)
	roleRepository            repository.RoleRepository            = repository.NewRoleRepository(db)
	businessTypeRepository    repository.BusinessTypeRepository    = repository.NewBusinessTypeRepository(db)
	paymentMethodRepository   repository.PaymentMethodRepository   = repository.NewPaymentMethodRepository(db)
	productCategoryRepository repository.ProductCategoryRepository = repository.NewProductCategoryRepository(db)
	productRepository         repository.ProductRepository         = repository.NewProductRepository(db)
	productVariantRepository  repository.ProductVariantRepository  = repository.NewProductVariantRepository(db)
	registrationRepository    repository.RegistrationRepository    = repository.NewRegistrationRepository(db)

	jwtService             service.JWTService             = service.NewJwtService()
	authService            service.AuthService            = service.NewAuthService(userRepository, userBusinessRepository)
	userService            service.UserService            = service.NewUserService(userRepository, validate)
	roleService            service.RoleService            = service.NewRoleService(roleRepository, validate)
	businessTypeService    service.BusinessTypeService    = service.NewBusinessTypeService(businessTypeRepository, validate)
	paymentMethodService   service.PaymentMethodService   = service.NewPaymentMethodService(paymentMethodRepository, validate)
	productCategoryService service.ProductCategoryService = service.NewProductCategoryService(productCategoryRepository, validate)
	registrationService    service.RegistrationService    = service.NewRegistrationService(registrationRepository, validate)
	productService         service.ProductService         = service.NewProductService(productRepository, productVariantRepository, validate)

	authController            controller.AuthController            = controller.NewAuthController(authService, jwtService)
	userController            controller.UserController            = controller.NewUserController(userService, jwtService)
	roleController            controller.RoleController            = controller.NewRoleController(roleService, jwtService)
	businessTypeController    controller.BusinessTypeController    = controller.NewBusinessTypeController(businessTypeService, jwtService)
	paymentMethodController   controller.PaymentMethodController   = controller.NewPaymentMethodController(paymentMethodService, jwtService)
	productCategoryController controller.ProductCategoryController = controller.NewProductCategoryController(productCategoryService, jwtService)
	registrationController    controller.RegistrationController    = controller.NewRegistrationController(registrationService)
	productController         controller.ProductController         = controller.NewProductController(productService, jwtService)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/business", authController.LoginBusiness)
		authRoutes.POST("", authController.Login)
	}

	registrationRoutes := r.Group("api/registration")
	{
		registrationRoutes.POST("", registrationController.Register)
	}

	userRoutes := r.Group("api/user")
	{
		userRoutes.POST("/create", userController.CreateUser)
		userRoutes.PUT("/update", userController.UpdateUser)
	}

	roleRoutes := r.Group("api/role")
	{
		roleRoutes.POST("", roleController.CreateRole)
		roleRoutes.PATCH("/:roleId", roleController.UpdateRole)
		roleRoutes.GET("", roleController.FindRoleAll)
		roleRoutes.GET("/:roleId", roleController.FindRoleById)
		roleRoutes.DELETE("/:roleId", roleController.DeleteRole)
	}

	businessTypeRoutes := r.Group("api/business-type", middleware.AuthorizeJWT(jwtService))
	{
		businessTypeRoutes.POST("/", businessTypeController.CreateBusinessType)
		businessTypeRoutes.PATCH("/:businessTypeId", businessTypeController.UpdateBusinessType)
		businessTypeRoutes.GET("/", businessTypeController.FindBusinessTypeAll)
		businessTypeRoutes.GET("/:businessTypeId", businessTypeController.FindBusinessTypeById)
		businessTypeRoutes.DELETE("/:businessTypeId", businessTypeController.DeleteBusinessType)
	}

	paymentMethodRoutes := r.Group("api/payment-method")
	{
		paymentMethodRoutes.POST("/", paymentMethodController.CreatePaymentMethod)
		paymentMethodRoutes.PATCH("/:paymentMethodId", paymentMethodController.UpdatePaymentMethod)
		paymentMethodRoutes.GET("/", paymentMethodController.FindPaymentMethodAll)
		paymentMethodRoutes.GET("/:paymentMethodId", paymentMethodController.FindPaymentMethodById)
		paymentMethodRoutes.DELETE("/:paymentMethodId", paymentMethodController.DeletePaymentMethod)
	}

	productCategoryRoutes := r.Group("api/product-category", middleware.AuthorizeJWT(jwtService))
	{
		productCategoryRoutes.POST("/", productCategoryController.Create)
		productCategoryRoutes.PATCH("/:id", productCategoryController.Update)
		productCategoryRoutes.GET("/", productCategoryController.FindAll)
		productCategoryRoutes.GET("/:id", productCategoryController.FindById)
		productCategoryRoutes.GET("/business/:business_id", productCategoryController.FindByBusinessId)
		productCategoryRoutes.DELETE("/:id", productCategoryController.Delete)
	}

	productRoutes := r.Group("api/product", middleware.AuthorizeJWT(jwtService))
	{
		productRoutes.POST("/", productController.Create)
		productRoutes.PATCH("/:id", productController.Update)
		productRoutes.GET("/", productController.FindAll)
		productRoutes.GET("/:id", productController.FindById)
		productRoutes.DELETE("/:id", productController.Delete)
	}

	return r
}
