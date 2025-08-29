package request

import "github.com/google/uuid"

type EmployeeRequest struct {
	Name        string    `json:"name"`
	BusinessId  uuid.UUID `json:"business_id"`
	RoleId      int       `json:"role_id" binding:"required"`
	Email       *string   `json:"email"`
	PhoneNumber *string   `json:"phone_number"`
	Password    *string   `json:"password"`
	PinCode     string    `json:"pin_code" binding:"required"`
}

type EmployeeUpdateRequest struct {
	Name   string `json:"name"`
	RoleId int    `json:"role_id"`
}

type PinLoginRequest struct {
	BusinessId  uuid.UUID `json:"business_id" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	PinCode     string    `json:"pin_code" binding:"required"`
}

type OpenShiftRequest struct {
	CashierId   uuid.UUID `json:"cashier_id"`
	BusinessId  uuid.UUID `json:"business_id"`
	TerminalId  uuid.UUID `json:"terminal_id" binding:"required"`
	OpeningCash float64   `json:"opening_cash"`
	Notes       *string   `json:"notes"`
}

type CloseShiftRequest struct {
	ClosingCash float64 `json:"closing_cash" binding:"required"`
	Notes       *string `json:"notes"`
}
