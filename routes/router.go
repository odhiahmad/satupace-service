package routes

import (
	"time"

	"run-sync/config"
	"run-sync/controller"
	"run-sync/helper"
	"run-sync/middleware"
	"run-sync/repository"
	"run-sync/service"
	ws "run-sync/websocket"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	emailHelper *helper.EmailHelper = helper.NewEmailHelper()
	redisClient *redis.Client       = config.SetupRedisClient()
	redisHelper *helper.RedisHelper = helper.NewRedisHelper(redisClient)
	validate    *validator.Validate = validator.New()
	db          *gorm.DB            = config.SetupDatabaseConnection()
	jwtService  service.JWTService  = service.NewJwtService()

	// Repositories
	userRepository     repository.UserRepository              = repository.NewUserRepository(db)
	runnerProfileRepo  repository.RunnerProfileRepository     = repository.NewRunnerProfileRepository(db)
	runGroupRepo       repository.RunGroupRepository          = repository.NewRunGroupRepository(db)
	runGroupMemberRepo repository.RunGroupMemberRepository    = repository.NewRunGroupMemberRepository(db)
	runActivityRepo    repository.RunActivityRepository       = repository.NewRunActivityRepository(db)
	directMatchRepo    repository.DirectMatchRepository       = repository.NewDirectMatchRepository(db)
	directChatRepo     repository.DirectChatMessageRepository = repository.NewDirectChatMessageRepository(db)
	groupChatRepo      repository.GroupChatMessageRepository  = repository.NewGroupChatMessageRepository(db)
	userPhotoRepo      repository.UserPhotoRepository         = repository.NewUserPhotoRepository(db)
	safetyLogRepo      repository.SafetyLogRepository         = repository.NewSafetyLogRepository(db)
	smartWatchRepo     repository.SmartWatchRepository        = repository.NewSmartWatchRepository(db)
	biometricRepo      repository.BiometricRepository         = repository.NewBiometricRepository(db)

	// Matching Engine
	matchingEngine service.MatchingEngine = service.NewMatchingEngine(runnerProfileRepo, directMatchRepo, runGroupRepo, safetyLogRepo)

	// Services
	userService          service.UserService           = service.NewUserService(userRepository)
	runnerProfileService service.RunnerProfileService  = service.NewRunnerProfileService(runnerProfileRepo, userRepository)
	runGroupService      service.RunGroupService       = service.NewRunGroupService(runGroupRepo, userRepository)
	runGroupMemberSvc    service.RunGroupMemberService = service.NewRunGroupMemberService(runGroupMemberRepo, userRepository, runGroupRepo, db)
	runActivitySvc       service.RunActivityService    = service.NewRunActivityService(runActivityRepo, userRepository, runnerProfileRepo)
	directMatchSvc       service.DirectMatchService    = service.NewDirectMatchService(directMatchRepo, userRepository, directChatRepo, runnerProfileRepo, matchingEngine, db)
	safetyLogSvc         service.SafetyLogService      = service.NewSafetyLogService(safetyLogRepo, userRepository, db)
	exploreSvc           service.ExploreService        = service.NewExploreService(runnerProfileRepo, runGroupRepo, directMatchRepo)
	smartWatchSvc        service.SmartWatchService     = service.NewSmartWatchService(smartWatchRepo, runActivityRepo)
	biometricSvc         service.BiometricService      = service.NewBiometricService(biometricRepo, userRepository, jwtService, redisHelper)

	// Controllers
	authController           controller.AuthController           = controller.NewAuthController(userService, jwtService, redisHelper, emailHelper)
	userController           controller.UserController           = controller.NewUserController(userService)
	runnerProfileController  controller.RunnerProfileController  = controller.NewRunnerProfileController(runnerProfileService)
	runGroupController       controller.RunGroupController       = controller.NewRunGroupController(runGroupService)
	runGroupMemberController controller.RunGroupMemberController = controller.NewRunGroupMemberController(runGroupMemberSvc)
	runActivityController    controller.RunActivityController    = controller.NewRunActivityController(runActivitySvc)
	directMatchController    controller.DirectMatchController    = controller.NewDirectMatchController(directMatchSvc)
	userPhotoController      controller.UserPhotoController      = controller.NewUserPhotoController(service.NewUserPhotoService(userPhotoRepo))
	safetyLogController      controller.SafetyLogController      = controller.NewSafetyLogController(safetyLogSvc)
	smartWatchController     controller.SmartWatchController     = controller.NewSmartWatchController(smartWatchSvc)
	exploreController        controller.ExploreController        = controller.NewExploreController(exploreSvc)
	biometricController      controller.BiometricController      = controller.NewBiometricController(biometricSvc)

	// WebSocket chat hub & controller
	chatHub          *ws.Hub                     = ws.NewHub()
	chatWSController controller.ChatWSController = controller.NewChatWSController(chatHub, directChatRepo, groupChatRepo, userRepository, jwtService)

	// WhatsApp controller
	whatsappController controller.WhatsAppController = controller.NewWhatsAppController(emailHelper, redisClient)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Start WebSocket hub in background
	go chatHub.Run()

	// Reusable middleware combos
	jwt := middleware.AuthorizeJWT(jwtService)
	profileReq := middleware.ProfileRequired(userRepository)

	// Auth routes (public)
	auth := r.Group("auth", middleware.RateLimit(redisHelper, 10, time.Minute))
	{
		auth.POST("/register", authController.Register)
		auth.POST("/verify", authController.VerifyOTP)
		auth.POST("/login", authController.Login)
		auth.POST("/resend-otp", authController.ResendOTP)

		// Biometric login (public - no JWT required)
		auth.POST("/biometric/login/start", biometricController.LoginStart)
		auth.POST("/biometric/login/finish", biometricController.LoginFinish)
	}

	// Biometric registration & management (JWT required)
	biometric := r.Group("biometric", jwt)
	{
		biometric.POST("/register/start", biometricController.RegisterStart)
		biometric.POST("/register/finish", biometricController.RegisterFinish)
		biometric.GET("/credentials", biometricController.GetCredentials)
		biometric.DELETE("/credentials/:id", biometricController.DeleteCredential)
	}

	// Public user routes
	users := r.Group("users", middleware.RateLimit(redisHelper, 20, time.Minute))
	{
		users.POST("", userController.Create)
		users.GET("", userController.FindAll)
		users.GET(":id", userController.FindById)
		users.PUT(":id", userController.Update)
		users.DELETE(":id", userController.Delete)
	}

	// Runner profile (JWT required, no profile required for create)
	profiles := r.Group("profiles", jwt)
	{
		profiles.POST("", runnerProfileController.CreateOrUpdate)
		profiles.GET("/me", runnerProfileController.FindByUserId)
		profiles.GET("/:id", runnerProfileController.FindById)
		profiles.PUT("/:id", runnerProfileController.Update)
		profiles.DELETE("/:id", runnerProfileController.Delete)
		profiles.GET("", runnerProfileController.FindAll)
	}

	// Run groups & activities
	runs := r.Group("runs")
	{
		// Groups
		runs.POST("/groups", jwt, runGroupController.Create)
		runs.GET("/groups", runGroupController.FindAll)
		runs.GET("/groups/:id", runGroupController.FindById)
		runs.PUT("/groups/:id", jwt, runGroupController.Update)
		runs.DELETE("/groups/:id", jwt, runGroupController.Delete)

		// Group members
		runs.POST("/groups/:id/join", jwt, profileReq, runGroupMemberController.JoinGroup)
		runs.GET("/groups/:id/members", runGroupMemberController.FindByGroupId)

		// Member management
		runs.PUT("/members/:id", jwt, runGroupMemberController.Update)
		runs.DELETE("/members/:id", jwt, runGroupMemberController.Delete)

		// Activities
		runs.POST("/activities", jwt, profileReq, runActivityController.Create)
		runs.GET("/activities/:id", runActivityController.FindById)
		runs.GET("/users/:userId/activities", runActivityController.FindByUserId)
	}

	// Explore / Discover (requires profile)
	explore := r.Group("explore", jwt, profileReq)
	{
		explore.GET("/runners", exploreController.FindNearbyRunners)
		explore.GET("/groups", exploreController.FindNearbyGroups)
	}

	// Direct matches (requires profile)
	dating := r.Group("match", jwt, profileReq)
	{
		dating.GET("/candidates", directMatchController.GetCandidates)
		dating.POST("", directMatchController.SendMatchRequest)
		dating.PATCH("/:id/accept", directMatchController.AcceptMatch)
		dating.PATCH("/:id/reject", directMatchController.RejectMatch)
		dating.GET("/:id", directMatchController.FindById)
		dating.GET("/me", directMatchController.FindUserMatches)
	}

	// WebSocket chat endpoints
	wsGroup := r.Group("ws")
	{
		wsGroup.GET("/direct/:matchId", chatWSController.HandleDirectChat)
		wsGroup.GET("/group/:groupId", chatWSController.HandleGroupChat)
	}

	// REST chat endpoints (history + delete)
	chats := r.Group("chats", jwt)
	{
		chats.GET("/direct/:matchId", chatWSController.GetDirectHistory)
		chats.GET("/group/:groupId", chatWSController.GetGroupHistory)
		chats.DELETE("/messages/:id", chatWSController.DeleteMessage)
	}

	// User photos & safety reports
	media := r.Group("media")
	{
		media.POST("/photos", jwt, userPhotoController.Create)
		media.GET("/photos/:id", userPhotoController.FindById)
		media.GET("/users/:userId/photos", userPhotoController.FindByUserId)

		media.POST("/safety", jwt, safetyLogController.ReportUser)
		media.GET("/safety/:id", safetyLogController.FindById)
	}

	// SmartWatch sync (requires auth + profile)
	watch := r.Group("smartwatch", jwt, profileReq)
	{
		watch.POST("/devices", smartWatchController.ConnectDevice)
		watch.GET("/devices", smartWatchController.GetDevices)
		watch.DELETE("/devices/:device_id", smartWatchController.DisconnectDevice)
		watch.GET("/devices/:device_id/stats", smartWatchController.GetDeviceStats)

		watch.POST("/sync", smartWatchController.SyncActivity)
		watch.POST("/sync/batch", smartWatchController.BatchSync)
		watch.GET("/sync/history", smartWatchController.GetSyncHistory)
	}

	// WhatsApp
	wa := r.Group("whatsapp", middleware.RateLimit(redisHelper, 10, time.Minute))
	{
		wa.POST("/register", whatsappController.Register)
		wa.POST("/verify", whatsappController.Verify)
	}

	return r
}
