package structs

type PingRequest struct {
	Version int
	Request PingRequestData
}

type PingResponse struct {
	Version  int
	Response PingResponseData
}

type PingRequestData struct {
	AgentVersion string
	AgentOS      string
	AgentID      string
	ClientID     string
}

type PingResponseData struct {
	Settings AgentSettings
	NewAgent AgentArchive
}

type AgentArchive struct {
	Zip []byte
}
