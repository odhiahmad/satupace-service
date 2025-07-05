package routes

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/config"
	"github.com/odhiahmad/kasirku-service/controller"
	"github.com/odhiahmad/kasirku-service/middleware"
	"github.com/odhiahmad/kasirku-service/repository"
	"github.com/odhiahmad/kasirku-service/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	resendFrom   = os.Getenv("RESEND_FROM")    // Contoh: "Kasirku <noreply@kasirku.com>"
	resendAPIKey = os.Getenv("RESEND_API_KEY") // Contoh: "re_live_xxx..."

	// Tambahkan log untuk debugging
	_ = fmt.Sprintf("📧 RESEND_FROM: %s", resendFrom)
	_ = fmt.Sprintf("🔐 RESEND_API_KEY: %s", resendAPIKey) // ⚠️ Hapus ini di production
	// Initialize the validator
	emailService service.EmailService = service.NewEmailService(
		os.Getenv("RESEND_FROM"),    // e.g. "Kasirku <noreply@kasirku.com>"
		os.Getenv("RESEND_API_KEY"), // API key dari Resend
	)
	redisClient               *redis.Client                        = config.SetupRedisClient()
	validate                  *validator.Validate                  = validator.New()
	db                        *gorm.DB                             = config.SetupDatabaseConnection()
	userRepository            repository.UserRepository            = repository.NewUserRepository(db)
	userBusinessRepository    repository.UserBusinessRepository    = repository.NewUserBusinessRepository(db)
	roleRepository            repository.RoleRepository            = repository.NewRoleRepository(db)
	businessTypeRepository    repository.BusinessTypeRepository    = repository.NewBusinessTypeRepository(db)
	paymentMethodRepository   repository.PaymentMethodRepository   = repository.NewPaymentMethodRepository(db)
	productCategoryRepository repository.ProductCategoryRepository = repository.NewProductCategoryRepository(db)
	productRepository         repository.ProductRepository         = repository.NewProductRepository(db)
	productPromoRepository    repository.ProductPromoRepository    = repository.NewProductPromoRepository(db)
	productVariantRepository  repository.ProductVariantRepository  = repository.NewProductVariantRepository(db)
	registrationRepository    repository.RegistrationRepository    = repository.NewRegistrationRepository(db)
	bundleRepository          repository.BundleRepository          = repository.NewBundleRepository(db)
	taxRepository             repository.TaxRepository             = repository.NewTaxRepository(db)
	discountRepository        repository.DiscountRepository        = repository.NewDiscountRepository(db)
	unitRepository            repository.UnitRepository            = repository.NewUnitRepository(db)
	transactionRepository     repository.TransactionRepository     = repository.NewTransactionRepository(db)
	promoRepository           repository.PromoRepository           = repository.NewPromoRepository(db)
	businessBranchRepository  repository.BusinessBranchRepository  = repository.NewBusinessBranchRepository(db)
	businessRepository        repository.BusinessRepository        = repository.NewBusinessRepository(db)
	membershipRepository      repository.MembershipRepository      = repository.NewMembershipRepository(db)

	jwtService             service.JWTService             = service.NewJwtService()
	authService            service.AuthService            = service.NewAuthService(userRepository, userBusinessRepository, jwtService)
	userService            service.UserService            = service.NewUserService(userRepository, validate)
	roleService            service.RoleService            = service.NewRoleService(roleRepository, validate)
	businessTypeService    service.BusinessTypeService    = service.NewBusinessTypeService(businessTypeRepository, validate)
	paymentMethodService   service.PaymentMethodService   = service.NewPaymentMethodService(paymentMethodRepository, validate)
	productCategoryService service.ProductCategoryService = service.NewProductCategoryService(productCategoryRepository, validate)
	registrationService    service.RegistrationService    = service.NewRegistrationService(registrationRepository, membershipRepository, emailService, validate)
	productService         service.ProductService         = service.NewProductService(productRepository, productPromoRepository, promoRepository, productVariantRepository, validate, redisClient)
	bundleService          service.BundleService          = service.NewBundleService(bundleRepository, validate)
	taxService             service.TaxService             = service.NewTaxService(taxRepository, validate)
	unitService            service.UnitService            = service.NewUnitService(unitRepository)
	transactionService     service.TransactionService     = service.NewTransactionService(db, transactionRepository, validate)
	promoService           service.PromoService           = service.NewPromoService(promoRepository, validate)
	discountService        service.DiscountService        = service.NewDiscountService(discountRepository, validate)
	businessBranchService  service.BusinessBranchService  = service.NewBusinessBranchService(businessBranchRepository, validate)
	businessService        service.BusinessService        = service.NewBusinessService(businessRepository, validate)
	productVariantService  service.ProductVariantService  = service.NewProductVariantService(productVariantRepository, productRepository, validate)

	authController            controller.AuthController            = controller.NewAuthController(authService, jwtService)
	userController            controller.UserController            = controller.NewUserController(userService, jwtService)
	roleController            controller.RoleController            = controller.NewRoleController(roleService, jwtService)
	businessTypeController    controller.BusinessTypeController    = controller.NewBusinessTypeController(businessTypeService, jwtService)
	paymentMethodController   controller.PaymentMethodController   = controller.NewPaymentMethodController(paymentMethodService, jwtService)
	productCategoryController controller.ProductCategoryController = controller.NewProductCategoryController(productCategoryService, jwtService)
	registrationController    controller.RegistrationController    = controller.NewRegistrationController(registrationService)
	productController         controller.ProductController         = controller.NewProductController(productService, jwtService)
	bundleController          controller.BundleController          = controller.NewBundleController(bundleService, jwtService)
	taxController             controller.TaxController             = controller.NewTaxController(taxService, jwtService)
	unitController            controller.UnitController            = controller.NewUnitController(unitService, jwtService)
	transactionController     controller.TransactionController     = controller.NewTransactionController(transactionService, jwtService)
	promoController           controller.PromoController           = controller.NewPromoController(promoService, jwtService)
	discountController        controller.DiscountController        = controller.NewDiscountController(discountService, jwtService)
	businessBranchController  controller.BusinessBranchController  = controller.NewBusinessBranchController(businessBranchService, jwtService)
	businessController        controller.BusinessController        = controller.NewBusinessController(businessService, jwtService)
	productVariantController  controller.ProductVariantController  = controller.NewProductVariantController(productVariantService, jwtService)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authRoutes := r.Group("auth")
	{
		authRoutes.POST("/business", authController.LoginBusiness)
		authRoutes.POST("/business/verify-email", authController.VerifyEmail)
		authRoutes.POST("", authController.Login)
	}

	registrationRoutes := r.Group("registration")
	{
		registrationRoutes.POST("", registrationController.Register)
	}

	userRoutes := r.Group("user")
	{
		userRoutes.POST("/create", userController.CreateUser)
		userRoutes.PUT("/update", userController.UpdateUser)
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

	productCategoryRoutes := r.Group("product-category", middleware.AuthorizeJWT(jwtService))
	{
		productCategoryRoutes.POST("", productCategoryController.Create)
		productCategoryRoutes.PATCH("/:id", productCategoryController.Update)
		productCategoryRoutes.GET("", productCategoryController.FindAll)
		productCategoryRoutes.GET("/:id", productCategoryController.FindById)
		productCategoryRoutes.GET("/business/:id", productCategoryController.FindByBusinessId)
		productCategoryRoutes.DELETE("/:id", productCategoryController.Delete)
	}

	productRoutes := r.Group("product", middleware.AuthorizeJWT(jwtService))
	{
		productRoutes.POST("", productController.Create)
		productRoutes.PATCH("/:id", productController.Update)
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
		bundleRoutes.PATCH("/:id", bundleController.Update)
		bundleRoutes.GET("/:id", bundleController.FindById)
		bundleRoutes.DELETE("/:id", bundleController.Delete)
		bundleRoutes.GET("", bundleController.FindWithPagination)
		bundleRoutes.PUT("/:id/active", bundleController.SetIsActive)
	}

	taxRoutes := r.Group("tax", middleware.AuthorizeJWT(jwtService))
	{
		taxRoutes.POST("", taxController.Create)
		taxRoutes.PATCH("/:id", taxController.Update)
		taxRoutes.GET("/:id", taxController.FindById)
		taxRoutes.DELETE("/:id", taxController.Delete)
		taxRoutes.GET("/business", taxController.FindWithPagination)
	}

	unitRoutes := r.Group("unit", middleware.AuthorizeJWT(jwtService))
	{
		unitRoutes.POST("", unitController.Create)
		unitRoutes.PATCH("/:id", unitController.Update)
		unitRoutes.GET("/:id", unitController.FindById)
		unitRoutes.DELETE("/:id", unitController.Delete)
		unitRoutes.GET("", unitController.FindWithPagination)
	}

	transactionRoutes := r.Group("transaction", middleware.AuthorizeJWT(jwtService))
	{
		transactionRoutes.POST("", transactionController.Create)
		transactionRoutes.PATCH("/:id/payment", transactionController.Payment)
		transactionRoutes.GET("/:id", transactionController.FindById)
		transactionRoutes.GET("", transactionController.FindWithPagination)
		transactionRoutes.PATCH("/:id/items", transactionController.AddOrUpdateItem)
	}

	promoRoutes := r.Group("promo", middleware.AuthorizeJWT(jwtService))
	{
		promoRoutes.POST("", promoController.Create)
		promoRoutes.PATCH("/:id", promoController.Update)
		promoRoutes.GET("/:id", promoController.FindById)
		promoRoutes.DELETE("/:id", promoController.Delete)
		promoRoutes.GET("/business", promoController.FindWithPagination)
		promoRoutes.PUT("/:id/active", promoController.SetIsActive)
	}

	discountRoutes := r.Group("discount", middleware.AuthorizeJWT(jwtService))
	{
		discountRoutes.POST("", discountController.Create)
		discountRoutes.PATCH("/:id", discountController.Update)
		discountRoutes.GET("/:id", discountController.FindById)
		discountRoutes.DELETE("/:id", discountController.Delete)
		discountRoutes.GET("/business", discountController.FindWithPagination)
		discountRoutes.PUT("/:id/active", discountController.SetIsActive)
	}

	businessBranchRoutes := r.Group("business-branch", middleware.AuthorizeJWT(jwtService))
	{
		businessBranchRoutes.POST("", businessBranchController.Create)
		businessBranchRoutes.PATCH("/:id", businessBranchController.Update)
		businessBranchRoutes.GET("/:id", businessBranchController.FindById)
		businessBranchRoutes.DELETE("/:id", businessBranchController.Delete)
		businessBranchRoutes.GET("", businessBranchController.FindWithPagination)
	}

	businessRoutes := r.Group("business", middleware.AuthorizeJWT(jwtService))
	{
		businessRoutes.POST("", businessController.Create)
		businessRoutes.PATCH("/:id", businessController.Update)
		businessRoutes.GET("/:id", businessController.FindById)
		businessRoutes.DELETE("/:id", businessController.Delete)
		businessRoutes.GET("", businessController.FindWithPagination)
	}

	return r
}
