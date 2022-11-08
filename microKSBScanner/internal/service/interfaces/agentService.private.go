package interfaces

import (
	"encoding/json"
	"net/http"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
)

func (ms *agentService) serverInfo(w http.ResponseWriter, body interface{}) {
	b, err := json.Marshal(successResp(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	ms.logg.Inf("serverInfo: " + string(b))
}

func (ms *agentService) clientError(w http.ResponseWriter, body interface{}) {
	b, err := json.Marshal(clientErrResp(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	ms.logg.Inf("clientError: " + string(b))
}

func (ms *agentService) serverError(w http.ResponseWriter, body interface{}) {
	b, err := json.Marshal(errorResp(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	ms.logg.Inf("serverError: " + string(b))
}

// type ResponseBody struct {
// 	StatusCode int         `json:"status_code"`
// 	Body       interface{} `json:"body"`
// }

// type ResponseError struct {
// 	StatusCode int         `json:"status_code"`
// 	Message    interface{} `json:"message"`
// }

func successResp(val interface{}) structs.ResponseBody {
	return structs.ResponseBody{
		StatusCode: 200,
		Body:       val,
	}
}

func errorResp(val interface{}) structs.ResponseBody {
	return structs.ResponseBody{
		StatusCode: 500,
		Message:    val,
	}
}

func clientErrResp(val interface{}) structs.ResponseBody {
	return structs.ResponseBody{
		StatusCode: 400,
		Message:    val,
	}
}
