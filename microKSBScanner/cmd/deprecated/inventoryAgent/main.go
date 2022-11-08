package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/domain"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/scanner"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/service"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/settings"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/updater"

	agent "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
)

var version = "0.0.0"

func main() {
	clientID, serverAddress, err := decodeFileName(os.Args[0])
	if err != nil {
		fmt.Printf("Некорректный файл, скачайте ещё раз.\n")
		return
	}

	logg, err := NewLogger("log.txt")
	if err != nil {
		fmt.Printf("Ошибка создания лог файла: %s\n", err.Error())
		return
	}

	err = run(logg, clientID, serverAddress)
	if err != nil {
		fmt.Printf("Ошибка: %s\n", err.Error())
		err := logg.Close()
		if err != nil {
			fmt.Printf("Ошибка: %s\n", err.Error())
			return
		}
		return
	}

	_ = logg.Close()
	_ = logg.Delete()

	fmt.Println("Успешно")
}

func run(logg *logger, clientID, serverAddress string) error {
	settRepo := settRepo{}
	dataRepo := dataRepo{}
	sensData := domain.SensitiveData{
		ClientID: clientID,
	}

	updr := updater.New(logg)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	cli := service.NewClient(httpClient)

	scan := scanner.New()

	sett, err := settings.New(24, 24, []string{serverAddress, serverAddress})
	if err != nil {
		return fmt.Errorf("settings.New err: %w", err)
	}

	ag, err := agent.New(version, logg, &sett, &settRepo, &dataRepo, &sensData, updr, cli, scan)
	if err != nil {
		return fmt.Errorf("agent.New err: %w", err)
	}

	err = ag.ScanAndSendInventory()
	if err != nil {
		return fmt.Errorf("Сканирование и отправка: %w", err)
	}

	return nil
}
