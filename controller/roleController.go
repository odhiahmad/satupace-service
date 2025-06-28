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
		response := helper.BuildErrorResponse(
			"Data tidak valid",
			"bad_request",
			"role",
			err.Error(),
			nil,
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.roleService.CreateRole(roleCreate)
	response := helper.BuildResponse(true, "Berhasil membuat role", nil)
	ctx.JSON(http.StatusCreated, response)
}

func (c *roleController) UpdateRole(ctx *gin.Context) {
	var roleUpdate request.RoleUpdate
	err := ctx.ShouldBind(&roleUpdate)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Data tidak valid",
			"bad_request",
			"role",
			err.Error(),
			nil,
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	roleId := ctx.Param("roleId")
	id, err := strconv.Atoi(roleId)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Parameter roleId tidak valid",
			"bad_request",
			"roleId",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	roleUpdate.Id = id

	c.roleService.UpdateRole(roleUpdate)
	response := helper.BuildResponse(true, "Berhasil mengubah role", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *roleController) FindRoleAll(ctx *gin.Context) {
	roleResponse := c.roleService.FindAll()
	response := helper.BuildResponse(true, "Data role berhasil diambil", roleResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *roleController) FindRoleById(ctx *gin.Context) {
	roleId := ctx.Param("roleId")
	id, err := strconv.Atoi(roleId)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Parameter roleId tidak valid",
			"bad_request",
			"roleId",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	roleResponse := c.roleService.FindById(id)
	response := helper.BuildResponse(true, "Data role berhasil ditemukan", roleResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *roleController) DeleteRole(ctx *gin.Context) {
	roleId := ctx.Param("roleId")
	id, err := strconv.Atoi(roleId)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Parameter roleId tidak valid",
			"bad_request",
			"roleId",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	c.roleService.Delete(id)
	response := helper.BuildResponse(true, "Berhasil menghapus role", nil)
	ctx.JSON(http.StatusOK, response)
}
