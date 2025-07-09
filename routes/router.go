package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/config"
	"github.com/odhiahmad/kasirku-service/controller"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/middleware"
	"github.com/odhiahmad/kasirku-service/repository"
	"github.com/odhiahmad/kasirku-service/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	emailService service.EmailService = service.NewEmailService(
		os.Getenv("RESEND_FROM"),    // e.g. "Kasirku <noreply@kasirku.com>"
		os.Getenv("RESEND_API_KEY"), // API key dari Resend
	)
	redisClient              *redis.Client                       = config.SetupRedisClient()
	redisHelper              *helper.RedisHelper                 = helper.NewRedisHelper(redisClient)
	validate                 *validator.Validate                 = validator.New()
	db                       *gorm.DB                            = config.SetupDatabaseConnection()
	userRepository           repository.UserRepository           = repository.NewUserRepository(db)
	userBusinessRepository   repository.UserBusinessRepository   = repository.NewUserBusinessRepository(db)
	roleRepository           repository.RoleRepository           = repository.NewRoleRepository(db)
	businessTypeRepository   repository.BusinessTypeRepository   = repository.NewBusinessTypeRepository(db)
	paymentMethodRepository  repository.PaymentMethodRepository  = repository.NewPaymentMethodRepository(db)
	categoryRepository       repository.CategoryRepository       = repository.NewCategoryRepository(db)
	productRepository        repository.ProductRepository        = repository.NewProductRepository(db)
	productPromoRepository   repository.ProductPromoRepository   = repository.NewProductPromoRepository(db)
	productVariantRepository repository.ProductVariantRepository = repository.NewProductVariantRepository(db)
	registrationRepository   repository.RegistrationRepository   = repository.NewRegistrationRepository(db)
	bundleRepository         repository.BundleRepository         = repository.NewBundleRepository(db)
	taxRepository            repository.TaxRepository            = repository.NewTaxRepository(db)
	discountRepository       repository.DiscountRepository       = repository.NewDiscountRepository(db)
	unitRepository           repository.UnitRepository           = repository.NewUnitRepository(db)
	transactionRepository    repository.TransactionRepository    = repository.NewTransactionRepository(db)
	promoRepository          repository.PromoRepository          = repository.NewPromoRepository(db)
	businessRepository       repository.BusinessRepository       = repository.NewBusinessRepository(db)
	membershipRepository     repository.MembershipRepository     = repository.NewMembershipRepository(db)
	brandRepository          repository.BrandRepository          = repository.NewBrandRepository(db)

	jwtService            service.JWTService            = service.NewJwtService()
	authService           service.AuthService           = service.NewAuthService(userRepository, userBusinessRepository, jwtService, redisHelper)
	userService           service.UserService           = service.NewUserService(userRepository, validate)
	userBusinessService   service.UserBusinessService   = service.NewUserBusinessService(userBusinessRepository, redisHelper)
	roleService           service.RoleService           = service.NewRoleService(roleRepository, validate)
	businessTypeService   service.BusinessTypeService   = service.NewBusinessTypeService(businessTypeRepository, validate)
	paymentMethodService  service.PaymentMethodService  = service.NewPaymentMethodService(paymentMethodRepository, validate)
	categoryService       service.CategoryService       = service.NewCategoryService(categoryRepository, validate)
	registrationService   service.RegistrationService   = service.NewRegistrationService(registrationRepository, membershipRepository, emailService, validate, redisHelper)
	productService        service.ProductService        = service.NewProductService(productRepository, productPromoRepository, promoRepository, productVariantRepository, validate, redisClient)
	bundleService         service.BundleService         = service.NewBundleService(bundleRepository, validate)
	taxService            service.TaxService            = service.NewTaxService(taxRepository, validate)
	unitService           service.UnitService           = service.NewUnitService(unitRepository)
	transactionService    service.TransactionService    = service.NewTransactionService(db, transactionRepository, validate)
	discountService       service.DiscountService       = service.NewDiscountService(discountRepository, validate)
	businessService       service.BusinessService       = service.NewBusinessService(businessRepository, validate)
	productVariantService service.ProductVariantService = service.NewProductVariantService(productVariantRepository, productRepository, validate)
	brandService          service.BrandService          = service.NewBrandService(brandRepository, validate)

	authController           controller.AuthController           = controller.NewAuthController(authService, jwtService)
	userController           controller.UserController           = controller.NewUserController(userService, jwtService)
	userBusinessController   controller.UserBusinessController   = controller.NewUserBusinessController(userBusinessService, jwtService)
	roleController           controller.RoleController           = controller.NewRoleController(roleService, jwtService)
	businessTypeController   controller.BusinessTypeController   = controller.NewBusinessTypeController(businessTypeService, jwtService)
	paymentMethodController  controller.PaymentMethodController  = controller.NewPaymentMethodController(paymentMethodService, jwtService)
	categoryController       controller.CategoryController       = controller.NewCategoryController(categoryService, jwtService)
	registrationController   controller.RegistrationController   = controller.NewRegistrationController(registrationService)
	productController        controller.ProductController        = controller.NewProductController(productService, jwtService)
	bundleController         controller.BundleController         = controller.NewBundleController(bundleService, jwtService)
	taxController            controller.TaxController            = controller.NewTaxController(taxService, jwtService)
	unitController           controller.UnitController           = controller.NewUnitController(unitService, jwtService)
	transactionController    controller.TransactionController    = controller.NewTransactionController(transactionService, jwtService)
	discountController       controller.DiscountController       = controller.NewDiscountController(discountService, jwtService)
	businessController       controller.BusinessController       = controller.NewBusinessController(businessService, jwtService)
	productVariantController controller.ProductVariantController = controller.NewProductVariantController(productVariantService, jwtService)
	brandController          controller.BrandController          = controller.NewBrandController(brandService, jwtService)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authRoutes := r.Group("auth")
	{
		authRoutes.POST("/business", authController.LoginBusiness)
		authRoutes.POST("/verify-otp", authController.VerifyOTP)
		authRoutes.POST("/retry-otp", authController.RetryOTP)
		authRoutes.POST("/registration", registrationController.Register)
	}

	userRoutes := r.Group("user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userBusinessController.FindById)
		userRoutes.PUT("/change-phone", userBusinessController.ChangePhone)
		userRoutes.PUT("/change-email", userBusinessController.ChangeEmail)
		userRoutes.PUT("/change-password", userBusinessController.ChangePassword)
		userRoutes.PUT("/business", businessController.Update)
	}

	roleRoutes := r.Group("role")
	{
		roleRoutes.POST("", roleController.CreateRole)
		roleRoutes.PATCH("/:roleId", roleController.UpdateRole)
		roleRoutes.GET("", roleController.FindRoleAll)
		roleRoutes.GET("/:roleId", roleController.FindRoleById)
		roleRoutes.DELETE("/:roleId", roleController.DeleteRole)
	}

	businessTypeRoutes := r.Group("business-type")
	{
		businessTypeRoutes.POST("", businessTypeController.CreateBusinessType)
		businessTypeRoutes.PATCH("/:id", businessTypeController.UpdateBusinessType)
		businessTypeRoutes.GET("", businessTypeController.FindBusinessTypeAll)
		businessTypeRoutes.GET("/:id", businessTypeController.FindBusinessTypeById)
		businessTypeRoutes.DELETE("/:id", businessTypeController.DeleteBusinessType)
	}

	paymentMethodRoutes := r.Group("payment-method")
	{
		paymentMethodRoutes.POST("", paymentMethodController.Create)
		paymentMethodRoutes.PATCH("/:id", paymentMethodController.Update)
		paymentMethodRoutes.GET("", paymentMethodController.FindAll)
		paymentMethodRoutes.GET("/:id", paymentMethodController.FindById)
		paymentMethodRoutes.DELETE("/:id", paymentMethodController.Delete)
	}

	libRoutes := r.Group("lib", middleware.AuthorizeJWT(jwtService))
	{
		libRoutes.POST("/category", categoryController.Create)
		libRoutes.POST("/brand", brandController.Create)
		libRoutes.POST("/tax", taxController.Create)
		libRoutes.POST("/unit", unitController.Create)
		libRoutes.POST("/discount", discountController.Create)

		libRoutes.PUT("/category/:id", categoryController.Update)
		libRoutes.PUT("/brand/:id", brandController.Update)
		libRoutes.PUT("/tax/:id", taxController.Update)
		libRoutes.PUT("/unit/:id", unitController.Update)
		libRoutes.PUT("/discount/:id", discountController.Update)

		libRoutes.GET("/category/:id", categoryController.FindById)
		libRoutes.GET("/brand/:id", brandController.FindById)
		libRoutes.GET("/tax/:id", taxController.FindById)
		libRoutes.GET("/unit/:id", unitController.FindById)
		libRoutes.GET("/discount/:id", discountController.FindById)

		libRoutes.GET("/category", categoryController.FindWithPagination)
		libRoutes.GET("/brand", brandController.FindWithPagination)
		libRoutes.GET("/tax", taxController.FindWithPagination)
		libRoutes.GET("/unit", unitController.FindWithPagination)
		libRoutes.GET("/discount", discountController.FindWithPagination)

		libRoutes.DELETE("/category/:id", categoryController.Delete)
		libRoutes.DELETE("/brand/:id", brandController.Delete)
		libRoutes.DELETE("/tax/:id", taxController.Delete)
		libRoutes.DELETE("/unit/:id", unitController.Delete)
		libRoutes.DELETE("/discount/:id", discountController.Delete)

		libRoutes.PUT("/discount/:id/active", discountController.SetIsActive)

	}

	productRoutes := r.Group("product", middleware.AuthorizeJWT(jwtService))
	{
		productRoutes.POST("", productController.Create)
		productRoutes.PUT("/:id", productController.Update)
		productRoutes.PUT("/image/:id", productController.UpdateImage)
		productRoutes.GET("/:id", productController.FindById)
		productRoutes.DELETE("/:id", productController.Delete)
		productRoutes.GET("", productController.FindWithPagination)
		productRoutes.POST("/:id/variant", productVariantController.Create)
		productRoutes.PATCH("/:id/variant", productVariantController.Update)
		productRoutes.DELETE("/variant/:id", productVariantController.Delete)
		productRoutes.DELETE("/variant/product/:productId", productVariantController.DeleteByProductId)
		productRoutes.GET("/variant/:id", productVariantController.FindById)
		productRoutes.GET("/variant/product/:productId", productVariantController.FindByProductId)
		productRoutes.PUT("/:id/active", productController.SetActive)
		productRoutes.PUT("/:id/available", productController.SetAvailable)
		productRoutes.PUT("/variant/:id/active", productVariantController.SetActive)
		productRoutes.PUT("/variant/:id/available", productVariantController.SetAvailable)
	}

	bundleRoutes := r.Group("bundle", middleware.AuthorizeJWT(jwtService))
	{
		bundleRoutes.POST("", bundleController.Create)
		bundleRoutes.PUT("/:id", bundleController.Update)
		bundleRoutes.GET("/:id", bundleController.FindById)
		bundleRoutes.DELETE("/:id", bundleController.Delete)
		bundleRoutes.GET("", bundleController.FindWithPagination)
		bundleRoutes.PUT("/:id/active", bundleController.SetIsActive)
		bundleRoutes.PUT("/:id/available", bundleController.SetIsAvailable)
	}

	transactionRoutes := r.Group("transaction", middleware.AuthorizeJWT(jwtService))
	{
		transactionRoutes.POST("", transactionController.Create)
		transactionRoutes.PUT("/:id/payment", transactionController.Payment)
		transactionRoutes.GET("/:id", transactionController.FindById)
		transactionRoutes.GET("", transactionController.FindWithPagination)
		transactionRoutes.PUT("/:id/items", transactionController.AddOrUpdateItem)
	}

	businessRoutes := r.Group("business", middleware.AuthorizeJWT(jwtService))
	{
		businessRoutes.POST("", businessController.Create)

		businessRoutes.GET("/:id", businessController.FindById)
		businessRoutes.DELETE("/:id", businessController.Delete)
		businessRoutes.GET("", businessController.FindWithPagination)
	}

	r.POST("/dev/reindex", func(ctx *gin.Context) {

		err := helper.CreateElasticProductIndex()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "Index berhasil dibuat"})
	})

	return r
}
