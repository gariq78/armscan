package winservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/judwhite/go-svc"
)

type service struct {
	nss     NameStartStoper
	ctx     context.Context
	exit    context.CancelFunc
	logFile *os.File
}

func (p *service) Context() context.Context {

	return p.ctx

}

func (p *service) Init(env svc.Environment) error {

	if env.IsWindowsService() {

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}

		logPath := filepath.Join(dir, p.nss.Name()+"service.log")

		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		p.logFile = f

		log.SetOutput(f)

	}

	return nil

}

func (p *service) Start() error {

	log.Printf("Starting...\n")

	go func(p *service) {

		err := p.nss.Start()
		if err != nil {

			log.Println(err.Error())
			p.exit()

			return

		}

		fmt.Printf("Started.\n")

	}(p)

	return nil

}

func (p *service) Stop() error {

	log.Printf("Stopping...\n")

	err := p.nss.Stop()
	if err != nil {
		log.Println(err.Error())
	}

	log.Printf("Stopped.\n")

	if p.logFile != nil {
		if closeErr := p.logFile.Close(); closeErr != nil {
			log.Printf("error closing '%s': %v\n", p.logFile.Name(), closeErr)
		}
	}

	return nil

}
