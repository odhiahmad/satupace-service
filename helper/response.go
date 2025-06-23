package helper

// Response is used for static shape of JSON return
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data"`
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
