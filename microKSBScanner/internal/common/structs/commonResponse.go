package structs

type ResponseBody struct {
	StatusCode int         `json:"status_code"`
	Body       interface{} `json:"body"`
	Message    interface{} `json:"message"`
}

// type ResponseError struct {
// 	StatusCode int         `json:"status_code"`
// 	Message    interface{} `json:"message"`
// }
