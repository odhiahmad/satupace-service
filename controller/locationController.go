package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type LocationController interface {
	GetProvinces(ctx *gin.Context)
	GetCities(ctx *gin.Context)
	GetDistricts(ctx *gin.Context)
	GetVillages(ctx *gin.Context)
}

type locationController struct {
	locationService service.LocationService
}

func NewLocationController(service service.LocationService) LocationController {
	return &locationController{locationService: service}
}

func (c *locationController) GetProvinces(ctx *gin.Context) {
	data, err := c.locationService.GetProvinces()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data provinsi", "ERROR", "server", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil data provinsi", data))
}

func (c *locationController) GetCities(ctx *gin.Context) {
	provinceIDStr := ctx.Query("province_id")
	provinceID, err := strconv.Atoi(provinceIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID provinsi tidak valid", "INVALID_ID", "province_id", err.Error(), nil))
		return
	}

	data, err := c.locationService.GetCitiesByProvince(provinceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data kota", "ERROR", "server", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil data kota", data))
}

func (c *locationController) GetDistricts(ctx *gin.Context) {
	cityIDStr := ctx.Query("city_id")
	cityID, err := strconv.Atoi(cityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID kota tidak valid", "INVALID_ID", "city_id", err.Error(), nil))
		return
	}

	data, err := c.locationService.GetDistrictsByCity(cityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data kecamatan", "ERROR", "server", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil data kecamatan", data))
}

func (c *locationController) GetVillages(ctx *gin.Context) {
	districtIDStr := ctx.Query("district_id")
	districtID, err := strconv.Atoi(districtIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID kecamatan tidak valid", "BAD_REQUEST", "district_id", err.Error(), nil))
		return
	}

	data, err := c.locationService.GetVillagesByDistrict(districtID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data kelurahan", "GET_FAILED", "", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil data kelurahan", data))
}
