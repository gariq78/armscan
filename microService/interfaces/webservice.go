package interfaces

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"ksb-dev.keysystems.local/intgrsrv/microService/domain"
	"ksb-dev.keysystems.local/intgrsrv/microService/usecases"
)

type Logger interface {
	Err(error, string, ...interface{})
	Inf(string, ...interface{})
}

type MicroServiceV1 struct {
	Interactor *usecases.Interactor
	Logger     Logger
	JWTParams  JWTParams
}

func (ms *MicroServiceV1) About(w http.ResponseWriter, r *http.Request) {
	type t struct {
		Name               string `json:"name"`
		Version            string `json:"version"`
		IntegrationName    string `json:"integration_name"`
		IntegrationVersion string `json:"integration_version"`
	}

	ab := ms.Interactor.About()

	var about t = t{
		Name:               ab.Name,
		Version:            ab.Version,
		IntegrationName:    ab.Source.Name,
		IntegrationVersion: ab.Source.Version,
	}

	ms.serverInfo(w, about)
}

func (ms *MicroServiceV1) StartAssets(w http.ResponseWriter, r *http.Request) {

	if ms.Interactor.StartAssets() {
		ms.serverInfo(w, nil)
	} else {
		ms.clientError(w, "Занято!") // сообщить что занят
	}
}

func (ms *MicroServiceV1) GetAssets(w http.ResponseWriter, r *http.Request) {

	d, err := ms.Interactor.GetAssets()

	if err != nil {
		if errors.Is(err, usecases.Busy) {
			ms.Logger.Inf("Занято!")
			ms.clientError(w, "Занято!") // сообщить что занят
			return
		}
		ms.serverError(w, fmt.Sprintf("Scan error: %s", err.Error()))
		return
	}

	ms.serverInfo(w, struct {
		Assets []domain.Asset `json:"assets"`
	}{
		Assets: d,
	})
}

const stampLayout = "2006-01-02T15:04:05Z"

func (ms *MicroServiceV1) Incidents(w http.ResponseWriter, r *http.Request) {
	var stamp time.Time
	var err error

	stampStr := r.FormValue("stamp")
	if stampStr != "" {
		stamp, err = time.Parse(stampLayout, stampStr)
		if err != nil {
			ms.clientError(w, fmt.Sprintf("stamp parse error: %s", err.Error()))
			return
		}
	}

	d, newStamp, err := ms.Interactor.Incidents(stamp)
	if err != nil {
		ms.serverError(w, fmt.Sprintf("Incidents error: %s", err.Error()))
		return
	}

	ms.serverInfo(w, struct {
		Incstamp  string            `json:"incstamp"`
		Incidents []domain.Incident `json:"incidents"`
	}{
		Incstamp:  newStamp.Format(stampLayout),
		Incidents: d,
	})
}

func (ms *MicroServiceV1) Settings(w http.ResponseWriter, r *http.Request) {
	sett := ms.Interactor.Settings()

	ms.serverInfo(w, sett)
}

func (ms *MicroServiceV1) SetSettings(w http.ResponseWriter, r *http.Request) {
	sett := r.PostFormValue("settings")

	err := ms.Interactor.SetSettings([]byte(sett))
	if err != nil {
		ms.serverError(w, fmt.Sprintf("interactor.SetSettings [%s]: %s", string(sett), err.Error()))
		return
	}

	ms.serverInfo(w, nil)
}
