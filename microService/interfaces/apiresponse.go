package interfaces

type ResponseBody struct {
	StatusCode int         `json:"status_code"`
	Body       interface{} `json:"body"`
}

type ResponseError struct {
	StatusCode int         `json:"status_code"`
	Message    interface{} `json:"message"`
}

func successResp(val interface{}) ResponseBody {
	return ResponseBody{
		StatusCode: 200,
		Body:       val,
	}
}

func errorResp(val interface{}) ResponseError {
	return ResponseError{
		StatusCode: 500,
		Message:    val,
	}
}

func clientErrResp(val interface{}) ResponseError {
	return ResponseError{
		StatusCode: 400,
		Message:    val,
	}
}
