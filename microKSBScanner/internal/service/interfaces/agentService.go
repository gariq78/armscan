package interfaces

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
	"net/http"
	"os"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
	service "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/service/usecases"
)

type agentService struct {
	logg       Logger
	interactor *service.KSBScannerSource
}

func NewAgentService(logg Logger, interactor *service.KSBScannerSource) (*agentService, error) {
	rv := &agentService{
		logg:       logg,
		interactor: interactor,
	}

	return rv, nil
}

func (as *agentService) AddData(w http.ResponseWriter, r *http.Request) {
	var p structs.AddDataRequset

	d := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := d.Decode(&p)

	if err != nil {
		as.clientError(w, fmt.Sprintf("AddData: AddDataRequest Decode error: %s", err.Error()))
		return
	}
	as.saveAssetsToJson(p.Request.Asset)

	as.logg.Inf("adddata request [%+v]", p)

	resp, err := as.interactor.AddData(p.Request)

	if err != nil {
		as.serverError(w, fmt.Sprintf("AddData: interactor error: %s", err.Error()))
		return
	}

	as.serverInfo(w, structs.PingResponse{
		Version: 1,
		Response: structs.PingResponseData{
			Settings: structs.AgentSettings{
				ServiceAddress:           resp.Settings.ServiceAddress,
				AdditionalServiceAddress: resp.Settings.AdditionalServiceAddress,
				PingPeriodHours:          resp.Settings.PingPeriodHours,
				MonitoringHours:          resp.Settings.MonitoringHours,
			},
			NewAgent: structs.AgentArchive{
				Zip: resp.NewAgent.Zip,
			},
		},
	})
}

func (as *agentService) Ping(w http.ResponseWriter, r *http.Request) {
	var p structs.PingRequest

	d := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := d.Decode(&p)
	if err != nil {
		as.clientError(w, fmt.Sprintf("Ping: PingRequest Decode error: %s", err.Error()))
		return
	}

	as.logg.Inf("ping request [%+v]", p)

	resp, err := as.interactor.Ping(structs.PingRequestData{
		AgentVersion: p.Request.AgentVersion,
		AgentID:      p.Request.AgentID,
		AgentOS:      p.Request.AgentOS,
	})
	if err != nil {
		as.serverError(w, fmt.Sprintf("Ping: interactor error: %s", err.Error()))
		return
	}

	as.serverInfo(w, structs.PingResponse{
		Version: 1,
		Response: structs.PingResponseData{
			Settings: structs.AgentSettings{
				ServiceAddress:           resp.Settings.ServiceAddress,
				AdditionalServiceAddress: resp.Settings.AdditionalServiceAddress,
				PingPeriodHours:          resp.Settings.PingPeriodHours,
				MonitoringHours:          resp.Settings.MonitoringHours,
			},
			NewAgent: structs.AgentArchive{
				Zip: resp.NewAgent.Zip,
			},
		},
	})
}

func (as *agentService) saveAssetsToJson(asset asset.Asset) {
	assetsFile := "./assets" + "/" + asset.HostName + ".json"

	saveStructToFile(asset, assetsFile)
}

func saveErrToFile(saveErr error) {
	fileName := "./errors.txt"

	isHas, err := checkFileToHas(fileName)
	if !isHas {
		emptyFile, errCreate := os.Create(fileName)
		if errCreate != nil {
			//errCollector(err, "func (f *FileWorker)SaveProductPositionsInCategory(positions marketStructs.PositionsInCategory)(err error)")
			return
		}
		defer emptyFile.Close()
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {

		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(saveErr.Error()); err != nil {
		panic(err)
	}
}

func saveStructToFile(someStruct interface{}, fileName string) (err error) {
	arrayBytes, errJsonMarshal := json.Marshal(someStruct)

	isHas, err := checkFileToHas(fileName)

	if !isHas {
		emptyFile, errCreate := os.Create(fileName)
		if errCreate != nil {
			//errCollector(err, "func (f *FileWorker)SaveProductPositionsInCategory(positions marketStructs.PositionsInCategory)(err error)")
			return
		}
		defer emptyFile.Close()
	}

	if errJsonMarshal != nil {
		//errCollector(errJsonMarshal, "func (f *FileWorker)SaveProductPositionsInCategory(positions marketStructs.PositionsInCategory)(err error)")
		return
	}
	err = ioutil.WriteFile(fileName, arrayBytes, 0644)
	return

}

func checkFileToHas(fileURL string) (isHas bool, err error) {
	_, err = os.Stat(fileURL)
	if err != nil {

		isHas = false
		return
	}
	isHas = true
	return
}
