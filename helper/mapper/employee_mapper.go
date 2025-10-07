package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
)

func MapEmployee(employee entity.Employee) *response.EmployeeResponse {

	var phone string
	if employee.PhoneNumber != nil {
		phone = *employee.PhoneNumber
	}

	return &response.EmployeeResponse{
		Id:          employee.Id,
		Name:        employee.Name,
		CreatedAt:   employee.CreatedAt,
		UpdatedAt:   employee.UpdatedAt,
		PhoneNumber: phone,
		Role:        MapRole(employee.Role),
		Business:    MapBusiness(employee.Business),
		IsActive:    employee.IsActive,
	}
}
