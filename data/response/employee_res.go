package response

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeResponse struct {
	Id          uuid.UUID         `json:"id"`
	Business    *BusinessResponse `json:"business"`
	PhoneNumber string            `json:"phone_number"`
	Name        string            `json:"name"`
	Role        *RoleResponse     `json:"role"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type EmployeeLoginResponse struct {
	Employee EmployeeResponse `json:"employee"`
}

type ShiftResponse struct {
	Id           uuid.UUID             `json:"id"`
	Business     *BusinessResponse     `json:"business"`
	Terminal     *TerminalResponse     `json:"terminal"`
	Cashier      *UserBusinessResponse `json:"cashier"`
	OpenedAt     time.Time             `json:"opened_at"`
	ClosedAt     *time.Time            `json:"closed_at,omitempty"`
	OpeningCash  float64               `json:"opening_cash"`
	ClosingCash  *float64              `json:"closing_cash,omitempty"`
	TotalSales   *float64              `json:"total_sales,omitempty"`
	TotalRefunds *float64              `json:"total_refunds,omitempty"`
	Status       string                `json:"status"`
	Notes        *string               `json:"notes,omitempty"`
}
