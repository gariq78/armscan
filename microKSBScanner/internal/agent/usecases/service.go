package usecases

import "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"

type Service interface {
	Ping(structs.AgentSettings, structs.PingRequestData) (structs.PingResponseData, error)
	SendData(structs.AgentSettings, structs.AddDataPacket) (structs.PingResponseData, error)
}

// type ServiceResponse struct {
// 	//Error    int // тут ошибки не должно быть, она должна возвращаться в результате методов, затем уже она должна заворачиваться в json перед ответом клиенту
// 	Settings Settings // как ты сюда анмаршал будешь делать? это же интерфейс
// 	NewAgent AgentArchive
// }

// type ServicePing struct {
// 	AgentVersion string
// 	AgentID      ID
// }
