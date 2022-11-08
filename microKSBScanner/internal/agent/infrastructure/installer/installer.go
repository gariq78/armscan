package installer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	gops "github.com/mitchellh/go-ps"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/domain"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
)

func Check(filename string, settRepo usecases.SettingsRepository, sensRepo usecases.SensitiveDataRepository) (structs.AgentSettings, domain.SensitiveData, error) {
	sett := structs.AgentSettings{}
	sens := domain.SensitiveData{}
	var err error

	// err = stopProcess(filename)
	// if err != nil {
	// 	return sett, sens, fmt.Errorf("stopProcess err: %w", err)
	// }

	err = removePrevVesion()
	if err != nil {
		return sett, sens, fmt.Errorf("removePrevVersion err: %w", err)
	}

	errSett := settRepo.Load(&sett)
	errSens := sensRepo.Load(&sens)
	clientID, serverAddr, errDeco := decodeFileName(filename)

	if errors.Is(errSett, os.ErrNotExist) {
		if errDeco != nil {
			return sett, sens, fmt.Errorf("settings file not found and parse filename err: %w", errDeco)
		}

		sett.ServiceAddress = serverAddr
		sett.PingPeriodHours = 1
		sett.MonitoringHours = 6

		err = settRepo.Save(sett)
		if err != nil {
			return sett, sens, fmt.Errorf("settRepo.Save err: %w", err)
		}
	} else if errSett != nil {
		return sett, sens, fmt.Errorf("settRepo.Load err: %w", errSett)
	}

	if errors.Is(errSens, os.ErrNotExist) {
		if errDeco != nil {
			return sett, sens, fmt.Errorf("id file not found and parse filename err: %w", errDeco)
		}

		sens.ClientID = clientID

		err = sensRepo.Save(sens)
		if err != nil {
			return sett, sens, fmt.Errorf("sensRepo.Save err: %w", err)
		}
	} else if errSens != nil {
		return sett, sens, fmt.Errorf("sensRepo.Load err: %w", errSens)
	}

	return sett, sens, nil
}

func stopProcess(name string) error {
	var procname = getOsDepProcName(name)
	ps, err := gops.Processes()
	if err != nil {
		return fmt.Errorf("processes err: %w", err)
	}
	var selfPid = os.Getpid()
	for _, p := range ps {
		var pid = p.Pid()
		if p.Executable() == procname && pid != selfPid { // сами себя не будем останавливать
			lav, err := os.FindProcess(pid)
			if err != nil {
				return fmt.Errorf("find process with pid [%d] err: %w", pid, err)
			}
			err = lav.Kill()
			if err != nil {
				return fmt.Errorf("kill process with pid [%d] err: %w", pid, err)
			}

			// Под виндой без этого говорит что не может перезаписать (при разархивации) тк процесс занят
			_, _ = lav.Wait()
			// Под линуксом здесь выходит по ошибке waitid: no child processes
			/*if e != nil {
				return e
			}*/
			log.Printf("process [%s] with pid [%d] stopped.", procname, pid)
			//break // закоментарил чтобы останавливались все процессы с таким именем, тк возможны разные порты
		}
	}
	return nil

}

func removePrevVesion() error {
	oldFiles, err := filepath.Glob("~*")
	if err != nil {
		return fmt.Errorf("glob err: %w", err)
	}

	if len(oldFiles) > 0 {
		for j := range oldFiles {
			var file = oldFiles[j]
			err = os.Remove(file)
			if err != nil {
				return fmt.Errorf("remove [%s] err: %w", file, err)
			} else {
				log.Printf("file [%s] was deleted", file)
			}
		}
	}

	return nil
}
