package handler

type Response struct {
	Data   interface{} `json:"data"`
	Status bool        `json:"status"`
	Errors []int64     `json:"errors"`
	err    error
}

// ResponseSuccess
func (r *Response) Success(data interface{}) *Response {
	return &Response{
		Data:   data,
		Status: true,
		Errors: []int64{0},
		err:    nil,
	}
}

// ResponseError
func (r *Response) Error(errorCode int64, error error) *Response {
	return &Response{
		Data:   error.Error(),
		Status: false,
		Errors: []int64{errorCode},
		err: error,
	}
}
