package usecases

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
	mdomain "ksb-dev.keysystems.local/intgrsrv/microService/domain"
)

type KSBScannerSource struct {
	about    mdomain.About
	settings Settings
	agentZip AgentZips

	logg    Logger
	settRep SettingsRepository
	client  *http.Client
	storer  DataRepository // TODO: выделить интерфейс для DataRepository
}

type OS string
type AgentZip struct {
	Version string
	Zip     []byte
}
type AgentZips map[OS]AgentZip

func NewKSBScannerSource(settRepo SettingsRepository, about mdomain.About, logg Logger, storer DataRepository, client *http.Client, zips AgentZips) (*KSBScannerSource, error) {
	var settings = Settings{
		Agent: structs.AgentSettings{
			ServiceAddress:           "http://localhost:3004/agent/api/v1",
			AdditionalServiceAddress: "http://localhost:3004/agent/api/v1",
			PingPeriodHours:          12,
			MonitoringHours:          6,
		},
	}
	err := settRepo.Load(&settings)
	if err != nil {
		logg.Inf("Can't load settings. Will be used default [%+v] values. err=[%s]", settings, err.Error())
		err = settRepo.Save(settings)
		if err != nil {
			return nil, fmt.Errorf("settingsRepo save err: %w", err)
		}
	}

	// storer, err := NewDataRepositoryClientID(dbHand, logg)
	// if err != nil {
	// 	return nil, fmt.Errorf("NewDataRepository err: %w", err)
	// }

	source := &KSBScannerSource{
		about:    about,
		settings: settings,
		agentZip: zips,

		logg:    logg,
		settRep: settRepo,
		client:  client,
		storer:  storer,
	}

	err = source.changeSettings(settings)
	if err != nil {
		return nil, fmt.Errorf("changeSettings(%+v): %w", settings, err)
	}

	return source, nil
}

func (s *KSBScannerSource) changeSettings(settings Settings) error {
	s.settings = settings

	return nil
}

func (s *KSBScannerSource) About() mdomain.About {
	return s.about
}

func (s *KSBScannerSource) Settings() interface{} {
	return plainSettings{
		AgentServiceAddress:           s.settings.Agent.ServiceAddress,
		AgentAdditionalServiceAddress: s.settings.Agent.AdditionalServiceAddress,
		AgentPingPeriodHours:          s.settings.Agent.PingPeriodHours,
		AgentMonitoringHours:          s.settings.Agent.MonitoringHours,
	}
}

func (s *KSBScannerSource) SetSettings(settings interface{}) error {
	var plain plainSettings
	err := json.Unmarshal(settings.([]byte), &plain)
	if err != nil {
		return fmt.Errorf("Unmarshal(%s): %w", string(settings.([]byte)), err)
	}

	s.logg.Inf("SetSettings: %#v", plain)

	sett := Settings{
		Agent: structs.AgentSettings{
			ServiceAddress:           plain.AgentServiceAddress,
			AdditionalServiceAddress: plain.AgentAdditionalServiceAddress,
			PingPeriodHours:          plain.AgentPingPeriodHours,
			MonitoringHours:          plain.AgentMonitoringHours,
		},
	}

	err = s.settRep.Save(sett)
	if err != nil {
		return fmt.Errorf("settRep.Save(%v): %w", sett, err)
	}

	err = s.changeSettings(sett)
	if err != nil {
		return fmt.Errorf("changeSettings=[%s]: %w", string(settings.([]byte)), err)
	}

	return nil
}

// TODO: Тут надо бы изменить интерфейс - добавить дату с которой возвращать данные по активам
func (s *KSBScannerSource) Assets() (assets []mdomain.Asset, err error) {
	datas, err := s.storer.Datas()

	if err != nil {
		return nil, fmt.Errorf("storer.Datas err: %w", err)
	}
	s.logg.Inf("source Assets count Datas %d", len(datas))

	assets = make([]mdomain.Asset, 0, len(datas))

	for _, d := range datas {

		cpus := make([]mdomain.CPUType, 0, len(d.CPU))
		for _, cpu := range d.CPU {
			cpus = append(cpus, mdomain.CPUType{
				Name: cpu.Name,
			})
		}

		hdds := make([]mdomain.HDDType, 0, len(d.HDD))
		for _, hdd := range d.HDD {
			hdds = append(hdds, mdomain.HDDType{
				Name:         hdd.Name,
				Manufacturer: hdd.Manufacturer,
				Size:         hdd.Size,
			})
		}

		softwares := make([]mdomain.SoftwareType, 0, len(d.Software))
		for _, software := range d.Software {
			softwares = append(softwares, mdomain.SoftwareType{
				Name:         software.Name,
				Manufacturer: software.Manufacturer,
				Version:      software.Version,
				Description:  software.Description,
			})
		}

		users := make([]mdomain.UsersType, 0, len(d.Users))
		for _, user := range d.Users {
			users = append(users, mdomain.UsersType{
				Login:       user.Login,
				FIO:         user.FIO,
				Description: user.Description,
				LoginDate:   user.LoginDate,
				Domain:      user.Domain,
			})
		}

		assets = append(assets, mdomain.Asset{
			ClientID:      d.ClientID,
			HostID:        d.HostID,
			HostName:      d.HostName,
			Domain:        d.Domain,
			IPAddress:     d.IPAddress,
			MACAddress:    d.MACAddress,
			Name:          d.Name,
			OSName:        d.OSName,
			OSVersion:     d.OSVersion,
			CPU:           cpus,
			SystemMemory:  d.SystemMemory,
			Video:         d.Video,
			HDD:           hdds,
			FactoryNumber: d.FactoryNumber,
			Motherboard:   d.Motherboard,
			Monitor:       d.Monitor,
			OD:            d.OD,
			Software:      softwares,
			Users:         users,
		})
	}

	return assets, nil
}

func (s *KSBScannerSource) Incidents(notBefore time.Time) (incidents []mdomain.Incident, newStamp time.Time, err error) {
	return nil, time.Time{}, fmt.Errorf("Not implemented Incidents")
}

// AGENT ---------------------------------------------------------------------------------------------

func (s *KSBScannerSource) Ping(p structs.PingRequestData) (structs.PingResponseData, error) {
	var rv = structs.PingResponseData{
		Settings: s.settings.Agent,
	}

	s.checkAgentVersion(p.AgentOS, p.AgentVersion, p.AgentID, &rv)

	return rv, nil
}

func (s *KSBScannerSource) AddData(a structs.AddDataPacket) (structs.PingResponseData, error) {
	var rv = structs.PingResponseData{
		Settings: s.settings.Agent,
	}

	err := s.storer.AddData(a)
	if err != nil {
		return rv, fmt.Errorf("storer.AddData err: %w", err)
	}

	s.checkAgentVersion(a.AgentOS, a.AgentVersion, a.AgentID, &rv)

	return rv, nil
}

func (s *KSBScannerSource) checkAgentVersion(agentos, agentversion, agentid string, rv *structs.PingResponseData) {
	if zip, ok := s.agentZip[OS(agentos)]; ok {
		if zip.Version != agentversion {
			rv.NewAgent = structs.AgentArchive{Zip: zip.Zip}
		}
	} else {
		s.logg.Inf("can not found agent zip for os %s pinged by agent %s %s", agentos, agentversion, agentid)
	}
}
