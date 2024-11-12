package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type MenuController interface {
	CreateMenu(ctx *gin.Context)
	UpdateMenu(ctx *gin.Context)
	FindMenuById(ctx *gin.Context)
	FindMenuAll(ctx *gin.Context)
	DeleteMenu(ctx *gin.Context)
}

type menuController struct {
	menuService service.MenuService
	jwtService  service.JWTService
}

func NewMenuController(menuService service.MenuService, jwtService service.JWTService) MenuController {
	return &menuController{
		menuService: menuService,
		jwtService:  jwtService,
	}
}

func (c *menuController) CreateMenu(ctx *gin.Context) {
	var menuCreate request.MenuCreate
	err := ctx.ShouldBind(&menuCreate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.menuService.CreateMenu(menuCreate)
	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *menuController) UpdateMenu(ctx *gin.Context) {
	menuUpdate := request.MenuUpdate{}
	err := ctx.ShouldBind(&menuUpdate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	menuId := ctx.Param("menuId")

	id, err := strconv.Atoi(menuId)
	helper.ErrorPanic(err)
	menuUpdate.Id = id

	c.menuService.UpdateMenu(menuUpdate)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *menuController) FindMenuAll(ctx *gin.Context) {
	menuResponse := c.menuService.FindAll()
	response := helper.BuildResponse(true, "!OK", menuResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *menuController) FindMenuById(ctx *gin.Context) {
	menuId := ctx.Param("menuId")
	id, err := strconv.Atoi(menuId)
	helper.ErrorPanic(err)

	menuResponse := c.menuService.FindById(id)

	response := helper.BuildResponse(true, "!OK", menuResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *menuController) DeleteMenu(ctx *gin.Context) {
	menuId := ctx.Param("menuId")
	id, err := strconv.Atoi(menuId)
	helper.ErrorPanic(err)

	c.menuService.Delete(id)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusOK, response)
}
