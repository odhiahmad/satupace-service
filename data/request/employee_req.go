package request

import "github.com/google/uuid"

type EmployeeRequest struct {
	Name        string    `json:"name"`
	BusinessId  uuid.UUID `json:"business_id"`
	RoleId      int       `json:"role_id" binding:"required"`
	PhoneNumber *string   `json:"phone_number"`
	Pin         string    `json:"pin" binding:"required"`
}

type EmployeeUpdateRequest struct {
	Name        *string `json:"name"`
	RoleId      *int    `json:"role_id"`
	Pin         *string `json:"pin"`
	PhoneNumber *string `json:"phone_number"`
	IsActive    *bool   `json:"is_active"`
}

type PinLoginRequest struct {
	EmployeeId uuid.UUID `json:"employee_id" binding:"required"`
	PinCode    string    `json:"pin_code" binding:"required"`
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
