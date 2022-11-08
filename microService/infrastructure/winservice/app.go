package winservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/judwhite/go-svc"
)

type NameStartStoper interface {
	Name() string
	Start() error
	Stop() error
}

type App struct {
	name    string
	service *service
}

func New(nss NameStartStoper) *App {
	ctx, cancel := context.WithCancel(context.Background())

	service := &service{
		nss:  nss,
		ctx:  ctx,
		exit: cancel,
	}

	return &App{
		name:    nss.Name(),
		service: service,
	}
}

// Run блокирует вызов, запускайте в горутине
func (t *App) Run(osArgs []string) {

	if len(osArgs) > 1 {
		err := t.serviceMode(osArgs[1])
		if err != nil {
			log.Fatalf("Service mode err: %s\n", err.Error())
			return
		}
		return
	}

	if err := svc.Run(t.service); err != nil {
		log.Fatal(err)
	}

}

func (t *App) serviceMode(mode string) error {

	if runtime.GOOS != "windows" {
		return fmt.Errorf("not implemented method")
	}

	switch mode {
	case "install":
		_, err := t.install()
		if err != nil {
			return fmt.Errorf("install err: %w", err)
		}
	case "uninstall":
		_, err := t.uninstall()
		if err != nil {
			return fmt.Errorf("uninstall err: %w", err)
		}
	case "start":
		_, err := t.start()
		if err != nil {
			return fmt.Errorf("start err: %w", err)
		}
	case "stop":
		_, err := t.stop()
		if err != nil {
			return fmt.Errorf("stop err: %w", err)
		}
	default:
		return fmt.Errorf("unknown command %q", mode)
	}

	return nil

}

type exitcode int
type exitcoder interface {
	ExitCode() int
}

func (t *App) install() (exitcode, error) {
	path, err := os.Executable()
	if err != nil {
		return 0, fmt.Errorf("executable err: %w", err)
	}

	code, err := t.execCmd("sc", "create", t.name, "binpath=", path, "start=", "auto", "DisplayName=", t.name)
	if err != nil {
		return code, fmt.Errorf("sc create err: %w", err)
	}

	code, err = t.execCmd("sc", "description", t.name, t.name+" service")
	if err != nil {
		return code, fmt.Errorf("sc description err: %w", err)
	}

	log.Printf(t.name + " service created")

	return t.start()
}

func (t *App) start() (exitcode, error) {
	code, err := t.execCmd("sc", "start", t.name)
	if err != nil {
		return code, fmt.Errorf("sc start err: %w", err)
	}

	log.Printf(t.name + " service started")

	return 0, nil
}

func (t *App) uninstall() (exitcode, error) {
	code, err := t.stop()
	if err != nil && code != 1062 { // если ошибка и ошибка не "сервис уже остановлен"
		return code, err
	}

	code, err = t.execCmd("sc", "delete", t.name)
	if err != nil {
		return code, fmt.Errorf("sc delete err: %w", err)
	}

	log.Printf(t.name + " service removed")

	return 0, nil
}

func (t *App) stop() (exitcode, error) {
	code, err := t.execCmd("sc", "stop", t.name)
	if err != nil {
		return code, fmt.Errorf("sc stop err: %w", err)
	}

	log.Printf(t.name + " service stopped")

	return 0, nil
}

func (t *App) execCmd(name string, args ...string) (exitcode, error) {

	cmd := exec.Command(name, args...)
	_, err := cmd.Output()

	if ee, ok := err.(exitcoder); ok {
		code := exitcode(ee.ExitCode())
		switch code {
		case 5: // нехватает прав
			return code, fmt.Errorf("please run as administrator.")
		case 1073: // уже установлена
			return code, fmt.Errorf(t.name + " service is already installed. For update first uninstall.")
		case 1056: // уже запущена
			return code, fmt.Errorf("already started.")
		case 1062: // уже остановлена
			return code, fmt.Errorf("already stopped.")
		case 1060: // не найдена
			return code, fmt.Errorf(t.name + " service not found. Please install first.")
		default:
			return code, err
		}
	}

	return 0, err

}
