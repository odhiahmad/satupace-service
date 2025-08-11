package response

import (
	"time"

	"github.com/google/uuid"
)

type EmployeePinLoginResponse struct {
	Token      string    `json:"token"`
	CashierId  uuid.UUID `json:"cashier_id"`
	Name       string    `json:"name"`
	RoleName   string    `json:"role_name"`
	BusinessId uuid.UUID `json:"business_id"`
	TerminalId uuid.UUID `json:"terminal_id"`
}

type ShiftResponse struct {
	Id           uuid.UUID  `json:"id"`
	BusinessId   uuid.UUID  `json:"business_id"`
	TerminalId   uuid.UUID  `json:"terminal_id"`
	CashierId    uuid.UUID  `json:"cashier_id"`
	CashierName  string     `json:"cashier_name"`
	OpenedAt     time.Time  `json:"opened_at"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	OpeningCash  float64    `json:"opening_cash"`
	ClosingCash  *float64   `json:"closing_cash,omitempty"`
	TotalSales   *float64   `json:"total_sales,omitempty"`
	TotalRefunds *float64   `json:"total_refunds,omitempty"`
	Status       string     `json:"status"`
	Notes        *string    `json:"notes,omitempty"`
}
