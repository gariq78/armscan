package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/infrastructure/installer"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/repository"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/scanner"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/service"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/updater"
	agent "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
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
	activeDataFileName := "data.bin"
	sensDataFileName := "id.bin"

	settFile := infrastructure.NewSettingsFile(settingsFileName, false)
	dataFile := infrastructure.NewSettingsFile(activeDataFileName, false) // тут название NewSettingsFile расходится с делом, надо бы переделать
	sensFile := infrastructure.NewSettingsFile(sensDataFileName, false)

	settRepo, err := repository.NewSettings(&settFile)
	if err != nil {
		return fmt.Errorf("repository.NewSettings err: %w", err)
	}

	dataRepo, err := repository.NewData(&dataFile)
	if err != nil {
		return fmt.Errorf("repository.NewData err: %w", err)
	}

	sensRepo, err := repository.NewSensitiveData(&sensFile)
	if err != nil {
		return fmt.Errorf("repository.NewSensitiveData err: %w", err)
	}

	sett, sensData, err := installer.Check(os.Args[0], settRepo, sensRepo)
	if err != nil {
		return fmt.Errorf("installer.Check err: %w", err)
	}

	updr := updater.New(logg, filepath.Base(os.Args[0]))

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	cli := service.NewClient(httpClient)

	scan := scanner.New()

	ag, err := agent.New(version, logg, sett, &settRepo, &dataRepo, &sensData, updr, cli, scan)
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
