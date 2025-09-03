package helper

import "loka-kasir/data/response"

type DetailedError struct {
	Code    string `json:"code"`
	Field   string `json:"field,omitempty"` // opsional
	Details string `json:"details"`
}

type ResponseError struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Error   *DetailedError `json:"error,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
}

// Response is used for static shape of JSON return
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data"`
}

type ResponsePagination struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors,omitempty"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}

type CursorPaginatedResponse struct {
	Limit      int    `json:"limit"`
	SortBy     string `json:"sort_by"`
	OrderBy    string `json:"order_by"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// EmptyObj is used when data should not be null in JSON
type EmptyObj struct{}

// BuildResponse creates a success response with data
func BuildResponse(status bool, message string, data interface{}) Response {
	if data == nil {
		data = EmptyObj{}
	}
	return Response{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func BuildResponsePagination(status bool, message string, data interface{}, pagination response.PaginatedResponse) ResponsePagination {
	if data == nil {
		data = EmptyObj{}
	}
	return ResponsePagination{
		Status:     status,
		Message:    message,
		Errors:     nil,
		Data:       data,
		Pagination: pagination,
	}
}

func BuildErrorResponse(message string, code string, field string, details string, data interface{}) ResponseError {
	if data == nil {
		data = EmptyObj{}
	}

	return ResponseError{
		Status:  false,
		Message: message,
		Error: &DetailedError{
			Code:    code,
			Field:   field,
			Details: details,
		},
		Data: data,
	}
}

func BuildResponseCursorPagination(status bool, message string, data interface{}, pagination response.CursorPaginatedResponse) ResponsePagination {
	if data == nil {
		data = EmptyObj{}
	}
	return ResponsePagination{
		Status:     status,
		Message:    message,
		Errors:     nil,
		Data:       data,
		Pagination: pagination,
	}
}
