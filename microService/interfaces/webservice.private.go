package interfaces

import (
	"encoding/json"
	"net/http"
)

func (ms *MicroServiceV1) serverInfo(w http.ResponseWriter, body interface{}) {
	b, err := json.Marshal(successResp(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

	var l = 1024
	if len(b) < l {
		l = len(b)
	}

	ms.Logger.Inf("serverInfo(%d/%d): %s", l, len(b), string(b[:l]))
}

func (ms *MicroServiceV1) clientError(w http.ResponseWriter, body interface{}) {
	b, err := json.Marshal(clientErrResp(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	ms.Logger.Inf("clientError: " + string(b))
}

func (ms *MicroServiceV1) serverError(w http.ResponseWriter, body interface{}) {
	b, err := json.Marshal(errorResp(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	ms.Logger.Inf("serverError: " + string(b))
}
