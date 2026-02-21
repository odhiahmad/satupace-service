package routes

import (
	"time"

	"run-sync/config"
	"run-sync/controller"
	"run-sync/helper"
	"run-sync/middleware"
	"run-sync/repository"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	emailHelper        *helper.EmailHelper                    = helper.NewEmailHelper()
	redisClient        *redis.Client                          = config.SetupRedisClient()
	redisHelper        *helper.RedisHelper                    = helper.NewRedisHelper(redisClient)
	validate           *validator.Validate                    = validator.New()
	db                 *gorm.DB                               = config.SetupDatabaseConnection()
	jwtService         service.JWTService                     = service.NewJwtService()
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

	// Application-level (non-business) user + run services/controllers
	userController           controller.UserController              = controller.NewUserController(service.NewUserService(userRepository))
	runnerProfileController  controller.RunnerProfileController     = controller.NewRunnerProfileController(service.NewRunnerProfileService(runnerProfileRepo, userRepository))
	runGroupController       controller.RunGroupController          = controller.NewRunGroupController(service.NewRunGroupService(runGroupRepo, userRepository))
	runGroupMemberController controller.RunGroupMemberController    = controller.NewRunGroupMemberController(service.NewRunGroupMemberService(runGroupMemberRepo, userRepository, runGroupRepo))
	runActivityController    controller.RunActivityController       = controller.NewRunActivityController(service.NewRunActivityService(runActivityRepo, userRepository))
	directMatchController    controller.DirectMatchController       = controller.NewDirectMatchController(service.NewDirectMatchService(directMatchRepo, userRepository))
	directChatController     controller.DirectChatMessageController = controller.NewDirectChatMessageController(service.NewDirectChatMessageService(directChatRepo, userRepository))
	groupChatController      controller.GroupChatMessageController  = controller.NewGroupChatMessageController(service.NewGroupChatMessageService(groupChatRepo, userRepository))
	userPhotoController      controller.UserPhotoController         = controller.NewUserPhotoController(service.NewUserPhotoService(userPhotoRepo))
	safetyLogController      controller.SafetyLogController         = controller.NewSafetyLogController(service.NewSafetyLogService(safetyLogRepo, userRepository))

	// WhatsApp controller (uses whatsmeow helper and Redis for mapping)
	whatsappController controller.WhatsAppController = controller.NewWhatsAppController(emailHelper, redisClient)
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public user routes (non-business users)
	users := r.Group("users", middleware.RateLimit(redisHelper, 20, time.Minute))
	{
		users.POST("", userController.Create)
		users.GET("", userController.FindAll)
		users.GET(":id", userController.FindById)
		users.PUT(":id", userController.Update)
		users.DELETE(":id", userController.Delete)
	}

	// Runner profile & run groups
	runs := r.Group("runs")
	{
		runs.POST("/groups", middleware.AuthorizeJWT(jwtService), runGroupController.Create)
		runs.GET("/groups", runGroupController.FindAll)
		runs.GET("/groups/:id", runGroupController.FindById)
		runs.PUT("/groups/:id", middleware.AuthorizeJWT(jwtService), runGroupController.Update)
		runs.DELETE("/groups/:id", middleware.AuthorizeJWT(jwtService), runGroupController.Delete)

		runs.POST("/groups/:id/join", middleware.AuthorizeJWT(jwtService), runGroupMemberController.JoinGroup)
		runs.GET("/groups/:groupId/members", runGroupMemberController.FindByGroupId)
		// Member management
		runs.PUT("/members/:id", middleware.AuthorizeJWT(jwtService), runGroupMemberController.Update)
		runs.DELETE("/members/:id", middleware.AuthorizeJWT(jwtService), runGroupMemberController.Delete)

		// Activities
		runs.POST("/activities", middleware.AuthorizeJWT(jwtService), runActivityController.Create)
		runs.GET("/activities/:id", runActivityController.FindById)
		runs.GET("/users/:userId/activities", runActivityController.FindByUserId)
	}

	// Direct matches & chats
	dating := r.Group("match")
	{
		dating.POST("", middleware.AuthorizeJWT(jwtService), directMatchController.Create)
		dating.PATCH("/:id", middleware.AuthorizeJWT(jwtService), directMatchController.Update)
		dating.GET("/:id", middleware.AuthorizeJWT(jwtService), directMatchController.FindById)
		dating.GET("/me", middleware.AuthorizeJWT(jwtService), directMatchController.FindUserMatches)
	}

	chats := r.Group("chats")
	{
		chats.POST("/direct", middleware.AuthorizeJWT(jwtService), directChatController.Create)
		chats.GET("/direct/:matchId", middleware.AuthorizeJWT(jwtService), directChatController.FindByMatchId)
		chats.POST("/group", middleware.AuthorizeJWT(jwtService), groupChatController.Create)
		chats.GET("/group/:groupId", middleware.AuthorizeJWT(jwtService), groupChatController.FindByGroupId)
	}

	// User photos & safety logs
	media := r.Group("media")
	{
		media.POST("/photos", middleware.AuthorizeJWT(jwtService), userPhotoController.Create)
		media.GET("/photos/:id", userPhotoController.FindById)
		media.GET("/users/:userId/photos", userPhotoController.FindByUserId)

		media.POST("/safety", middleware.AuthorizeJWT(jwtService), safetyLogController.Create)
		media.GET("/safety/:id", safetyLogController.FindById)
	}

	// WhatsApp registration and mapping
	wa := r.Group("whatsapp", middleware.RateLimit(redisHelper, 10, time.Minute))
	{
		wa.POST("/register", whatsappController.Register)
		wa.POST("/verify", whatsappController.Verify)
	}

	return r
}
