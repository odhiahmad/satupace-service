package helper

import "github.com/odhiahmad/kasirku-service/data/response"

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

// BuildErrorResponse creates an error response with error detail
func BuildErrorResponse(message string, err interface{}, data interface{}) Response {
	var errDetail interface{}

	switch e := err.(type) {
	case string:
		errDetail = []string{e}
	case error:
		errDetail = []string{e.Error()}
	case []string:
		errDetail = e
	default:
		errDetail = []string{"unexpected error format"}
	}

	if data == nil {
		data = EmptyObj{}
	}

	return Response{
		Status:  false,
		Message: message,
		Errors:  errDetail,
		Data:    data,
	}
}
