package mapper

import (
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapShift(shift *entity.Shift) *response.ShiftResponse {
	if shift == nil || shift.Id == uuid.Nil {
		return nil
	}

	return &response.ShiftResponse{
		Id:           shift.Id,
		Business:     MapBusiness(shift.Business),
		Terminal:     MapTerminal(shift.Terminal),
		Cashier:      MapUserBusiness(*shift.Cashier),
		OpenedAt:     shift.OpenedAt,
		ClosedAt:     shift.ClosedAt,
		OpeningCash:  shift.OpeningCash,
		ClosingCash:  shift.ClosingCash,
		TotalSales:   shift.TotalSales,
		TotalRefunds: shift.TotalRefunds,
		Status:       shift.Status,
		Notes:        shift.Notes,
	}
}
