package service

import (
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
)

// type PingRequest struct {
// 	Version int
// 	Ping    agent.ServicePing
// }

type PingResponse struct {
	StatusCode int `json:"status_code"`
	Body       struct {
		Version  int
		Response structs.PingResponseData
	} `json:"body"`
	Message string `json:"message"`
}
