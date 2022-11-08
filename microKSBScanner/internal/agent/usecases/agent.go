package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/additionInfo"
	"os"
	"runtime"
	"time"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/domain"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

type Agent struct {
	id       ID
	sensData domain.SensitiveData
	version  string
	data     asset.Asset
	settings structs.AgentSettings

	monitoringTimer       *time.Ticker
	monitoringTimerStoper chan struct{}

	pingingTimer       *time.Ticker
	pingingTimerStoper chan struct{}

	logger       Logger
	settingsRepo SettingsRepository
	dataRepo     DataRepository

	updater Updater
	service Service
	scanner Scanner

	pingReqData structs.PingRequestData
}

func New(
	version string,
	logger Logger,
	settings structs.AgentSettings,
	settingsRepo SettingsRepository,
	dataRepo DataRepository,
	sensData *domain.SensitiveData,
	updater Updater,
	service Service,
	scanner Scanner) (*Agent, error) {

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if settingsRepo == nil {
		return nil, fmt.Errorf("settingsRepo is nil")
	}

	if dataRepo == nil {
		return nil, fmt.Errorf("dataRepo is nil")
	}

	if service == nil {
		return nil, fmt.Errorf("service is nil")
	}

	if scanner == nil {
		return nil, fmt.Errorf("scanner is nil")
	}

	if updater == nil {
		return nil, fmt.Errorf("updater is nil")
	}

	id, err := scanner.ID()
	if err != nil {
		return nil, fmt.Errorf("scanner.ID err: %w", err)
	}

	logger.Inf("agentid %s", id)
	// это наверное нужно вынести наружу, так как там есть и репозиторий и имплементация настроек,
	// а сюда уже будет приходить нужные настройки (либо из параметров либо из сохранённого файла)
	// 2020.05.19: верно глаголишь - см. ниже работу с sensRepo
	/*
		err := repository.LoadSettings(settings)
		if err != nil {
			return nil, fmt.Errorf("repository.LoadSettings error: %w", err)
		}
	*/

	sensD := domain.SensitiveData{}

	if sensData != nil {
		sensD = *sensData
		logger.Inf("clientid %s", sensD.ClientID)
	}

	return &Agent{
		id:                    id,
		sensData:              sensD,
		version:               version,
		settings:              settings,
		monitoringTimerStoper: make(chan struct{}),
		pingingTimerStoper:    make(chan struct{}),
		logger:                logger,
		settingsRepo:          settingsRepo,
		dataRepo:              dataRepo,
		updater:               updater,
		service:               service,
		scanner:               scanner,
		pingReqData: structs.PingRequestData{
			AgentVersion: version,
			AgentOS:      runtime.GOOS,
			AgentID:      string(id),
			ClientID:     sensD.ClientID,
		},
	}, nil
}

// Стартовать должен в любом случае и должен мочь достучаться до сервера чтобы получить настройки либо обновление
func (agent *Agent) Start() error {
	agent.logger.Inf("starting %s version", agent.version)

	err := agent.loadSettings()
	if err != nil {
		return err
	}

	data, err := agent.loadData()

	if err != nil {
		agent.logger.Err(err, "Start: loadData err")
	}

	agent.data = data

	err = agent.scanAndSend()
	if err != nil {
		agent.logger.Err(err, "Start: scanAndSend err")
	}

	agent.startMonitoring(agent.settings.GetMonitoringDuration())
	agent.startPinging(agent.settings.GetPingDuration())

	return nil
}

func (agent *Agent) Stop() {
	if agent.monitoringTimer != nil {
		agent.monitoringTimer.Stop()
		agent.monitoringTimerStoper <- struct{}{}
	}
	if agent.pingingTimer != nil {
		agent.pingingTimer.Stop()
		agent.pingingTimerStoper <- struct{}{}
	}

	agent.logger.Inf("agent timers stoped")
}

func (agent *Agent) startMonitoring(duration time.Duration) {

	agent.monitoringTimer = time.NewTicker(duration)

	go func() {
		agent.logger.Inf("enter to restartMonitoring gorutine")
	f:
		for {
			select {
			case <-agent.monitoringTimer.C:
				agent.logger.Inf("start scanAndSend")
				err := agent.scanAndSend()
				if err != nil {
					agent.logger.Err(err, "ScanAndSend")
				}
			case <-agent.monitoringTimerStoper:
				agent.logger.Inf("monitoringTimerStoper")
				break f
			}
		}
		agent.logger.Inf("exiting from restartMonitoring gorutine")
	}()

}

func (agent *Agent) restartMonitoring(duration time.Duration) {

	agent.logger.Inf("restart monitoring %v", duration)

	agent.monitoringTimer.Reset(duration)

	agent.logger.Inf("monitoring restarted")

}

func (agent *Agent) startPinging(duration time.Duration) {

	agent.pingingTimer = time.NewTicker(duration)

	go func() {
		agent.logger.Inf("enter to restartPinging gorutine")
	f:
		for {
			select {
			case <-agent.pingingTimer.C:
				agent.logger.Inf("start ping")
				err := agent.ping()
				if err != nil {
					agent.logger.Err(err, "agent.ping")
				}
			case <-agent.pingingTimerStoper:
				agent.logger.Inf("pingingTimerStoper")
				break f
			}
		}
		agent.logger.Inf("exiting from restartPinging gorutine")
	}()

}

func (agent *Agent) restartPinging(duration time.Duration) {

	agent.logger.Inf("restart pinging %v", duration)

	agent.pingingTimer.Reset(duration)

	agent.logger.Inf("pinging restarted")

}

func (agent *Agent) ScanAndSendInventory() error {
	data, err := agent.scanner.Scan()
	if err != nil {
		return err
	}

	agent.logger.Inf("Scanned")

	if agent.data.Changed(data) {
		agent.logger.Inf("data changed need send to server")
		data.ClientID = agent.sensData.ClientID
		_, err := agent.service.SendData(agent.settings, structs.AddDataPacket{
			PingRequestData: agent.pingReqData,
			ScanTime:        time.Now().UTC(),
			Asset:           data,
		})
		if err != nil {
			// TODO: вызывать смену адреса сервиса и повторную отправку
			return fmt.Errorf("service.SendData err: %w", err)
		}
		agent.logger.Inf("data sended to server")
	} else {
		agent.logger.Inf("data not changed")
	}

	return nil
}

func (agent *Agent) scanAndSend() error {
	data, err := agent.scanner.Scan()
	if err != nil {
		return err
	}

	data.Software = additionInfo.AddInfo(data.Software)

	agent.logger.Inf("Scanned")

	if agent.data.Changed(data) {
		agent.logger.Inf("data changed need send to server")
		data.ClientID = agent.sensData.ClientID
		resp, err := agent.service.SendData(agent.settings, structs.AddDataPacket{
			PingRequestData: agent.pingReqData,
			ScanTime:        time.Now().UTC(),
			Asset:           data,
		})
		if err != nil {
			// TODO: вызывать смену адреса сервиса и повторную отправку
			return fmt.Errorf("service.SendData err: %w", err)
		}
		agent.logger.Inf("data sended to server")

		agent.data = data

		err = agent.saveData(data)
		if err != nil {
			return err
		}

		err = agent.applySettings(resp.Settings)
		if err != nil {
			return err
		}

		err = agent.update(resp.NewAgent)
		if err != nil {
			return err
		}
	} else {
		agent.logger.Inf("data not changed")
	}

	return nil
}

func (agent *Agent) ping() error {
	agent.logger.Inf("ping server")
	resp, err := agent.service.Ping(agent.settings, agent.pingReqData)
	if err != nil {
		// TODO: вызывать смену адреса сервиса и повторную отправку только если неотвечает
		return fmt.Errorf("service.Ping err: %w", err)
	}
	agent.logger.Inf("pinged server")

	err = agent.applySettings(resp.Settings)
	if err != nil {
		return fmt.Errorf("agent.applySettings err: %w", err)
	}

	err = agent.update(resp.NewAgent)
	if err != nil {
		return fmt.Errorf("agent.update err: %w", err)
	}

	return nil
}

func (agent *Agent) applySettings(n structs.AgentSettings) error {
	if agent.settings.NeedUpdate(n) {
		if agent.settings.GetMonitoringDuration() != n.GetMonitoringDuration() {
			agent.restartMonitoring(n.GetMonitoringDuration())
		}

		if agent.settings.GetPingDuration() != n.GetPingDuration() {
			agent.restartPinging(n.GetPingDuration())
		}

		jso, _ := json.Marshal(agent.settings)
		agent.logger.Inf("applySettings: old %+v", string(jso))
		agent.settings = n

		jsn, _ := json.Marshal(n)
		agent.logger.Inf("applySettings: new %+v", string(jsn))

		err := agent.saveSettings()
		if err != nil {
			return err
		}
	} else {
		// js, err := json.Marshal(n)
		// if err != nil {
		// 	agent.logger.Err(err, "debug marshal settings: %w")
		// }
		agent.logger.Inf("applySettings: settings not changed")
	}

	return nil
}

func (agent *Agent) loadSettings() error {
	sett := agent.settings
	err := agent.settingsRepo.Load(&sett)
	if err != nil {
		return err
	}

	js, _ := json.Marshal(sett)
	agent.logger.Inf("loadSettings: %+v", string(js))

	return agent.applySettings(sett)
}

func (agent *Agent) saveSettings() error {
	return agent.settingsRepo.Save(agent.settings)
}

func (agent *Agent) loadData() (asset.Asset, error) {
	data, err := agent.dataRepo.LoadData()

	// если первый старт или файл не найден, то начинаем сначала
	if errors.Is(err, os.ErrNotExist) {
		return asset.Asset{}, nil
	}

	return data, err
}

func (agent *Agent) saveData(d asset.Asset) error {
	return agent.dataRepo.SaveData(d)
}

func (agent *Agent) update(n structs.AgentArchive) error {
	err := agent.updater.Update(n)
	if err != nil {
		return err
	}

	return nil
}
