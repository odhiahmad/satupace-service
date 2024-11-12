package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type RoleController interface {
	CreateRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	FindRoleById(ctx *gin.Context)
	FindRoleAll(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
}

type roleController struct {
	roleService service.RoleService
	jwtService  service.JWTService
}

func NewRoleController(roleService service.RoleService, jwtService service.JWTService) RoleController {
	return &roleController{
		roleService: roleService,
		jwtService:  jwtService,
	}
}

func (c *roleController) CreateRole(ctx *gin.Context) {
	var roleCreate request.RoleCreate
	err := ctx.ShouldBind(&roleCreate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.roleService.CreateRole(roleCreate)
	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *roleController) UpdateRole(ctx *gin.Context) {
	roleUpdate := request.RoleUpdate{}
	err := ctx.ShouldBind(&roleUpdate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	roleId := ctx.Param("roleId")

	id, err := strconv.Atoi(roleId)
	helper.ErrorPanic(err)
	roleUpdate.Id = id

	c.roleService.UpdateRole(roleUpdate)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *roleController) FindRoleAll(ctx *gin.Context) {
	roleResponse := c.roleService.FindAll()
	response := helper.BuildResponse(true, "!OK", roleResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *roleController) FindRoleById(ctx *gin.Context) {
	roleId := ctx.Param("roleId")
	id, err := strconv.Atoi(roleId)
	helper.ErrorPanic(err)

	roleResponse := c.roleService.FindById(id)

	response := helper.BuildResponse(true, "!OK", roleResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *roleController) DeleteRole(ctx *gin.Context) {
	roleId := ctx.Param("roleId")
	id, err := strconv.Atoi(roleId)
	helper.ErrorPanic(err)

	c.roleService.Delete(id)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusOK, response)
}
