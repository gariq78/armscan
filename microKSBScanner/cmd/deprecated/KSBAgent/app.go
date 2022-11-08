package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/repository"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/scanner"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/service"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/updater"
	agent "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
	"ksb-dev.keysystems.local/intgrsrv/microService/infrastructure"
	"ksb-dev.keysystems.local/intgrsrv/microService/infrastructure/winservice"
)

type application struct {
	appversion string
	appname    string
	ag         *agent.Agent
}

var _ winservice.NameStartStoper = &application{}

func (t *application) Name() string {

	return t.appname

}

func (t *application) Stop() error {

	if t.ag != nil {
		t.ag.Stop()
	}

	return nil

}

func (t *application) Start() error {
	logg := infrastructure.NewMyUglyLogger("")
	defer logg.Close() // выглядит непонятно, ладно реализации там внутри нет и всё равно что нет вызова, но надо бы переделать

	settingsFileName := "settings.bin"
	activeDataFileName := "data.json"

	settFile := infrastructure.NewSettingsFile(settingsFileName, false)
	dataFile := infrastructure.NewSettingsFile(activeDataFileName, true) // тут название NewSettingsFile расходится с делом, надо бы переделать

	settRepo, err := repository.NewSettings(&settFile)
	if err != nil {
		return fmt.Errorf("repository.NewSettings err: %w", err)
	}

	dataRepo, err := repository.NewData(&dataFile)
	if err != nil {
		return fmt.Errorf("repository.NewData err: %w", err)
	}

	updr := updater.New(logg)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	cli := service.NewClient(httpClient)

	scan := scanner.New()

	sett := structs.AgentSettings{}

	err = settRepo.Load(&sett)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file %q not found", settingsFileName)
		} else {
			return fmt.Errorf("settRepo.Load err: %w", err)
		}
	}

	ag, err := agent.New(version, logg, sett, &settRepo, &dataRepo, nil, updr, cli, scan)
	if err != nil {
		return fmt.Errorf("agent.New err: %w", err)
	}

	err = ag.Start()
	if err != nil {
		return fmt.Errorf("agent.Start err: %w", err)
	}

	t.ag = ag

	return nil

}
