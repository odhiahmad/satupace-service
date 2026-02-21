package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExploreController interface {
	FindNearbyRunners(ctx *gin.Context)
	FindNearbyGroups(ctx *gin.Context)
}

type exploreController struct {
	service service.ExploreService
}

func NewExploreController(s service.ExploreService) ExploreController {
	return &exploreController{service: s}
}

func (c *exploreController) FindNearbyRunners(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.ExploreRunnersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		res := helper.BuildErrorResponse("Parameter tidak valid", "INVALID_PARAMS", "query", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	results, err := c.service.FindNearbyRunners(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mencari runner", "EXPLORE_FAILED", "query", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data runner terdekat", results)
	ctx.JSON(http.StatusOK, response)
}

func (c *exploreController) FindNearbyGroups(ctx *gin.Context) {
	var req request.ExploreGroupsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		res := helper.BuildErrorResponse("Parameter tidak valid", "INVALID_PARAMS", "query", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	results, err := c.service.FindNearbyGroups(req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mencari grup", "EXPLORE_FAILED", "query", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data grup terdekat", results)
	ctx.JSON(http.StatusOK, response)
}
