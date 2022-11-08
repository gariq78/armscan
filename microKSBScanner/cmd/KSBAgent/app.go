package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/service/infrastructure"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/service/infrastructure/webapp"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/service/interfaces"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/service/usecases"

	microDomain "ksb-dev.keysystems.local/intgrsrv/microService/domain"
	microInfrastructure "ksb-dev.keysystems.local/intgrsrv/microService/infrastructure"
	microInterfaces "ksb-dev.keysystems.local/intgrsrv/microService/interfaces"
	microUsecases "ksb-dev.keysystems.local/intgrsrv/microService/usecases"

	"ksb-dev.keysystems.local/intgrsrv/microService/infrastructure/winservice"
)

type application struct {
	appname    string
	appversion string

	intgrSrv *webapp.WebApp
	agentSrv *webapp.WebApp

	dbhand *infrastructure.SqliteHandler

	exit chan struct{}
}

var _ winservice.NameStartStoper = &application{}

func (t *application) Name() string {

	return t.appname

}

func (t *application) Stop() error {

	err := t.dbhand.Close()
	if err != nil {
		log.Printf("dbhand close err: %s", err.Error())
	}

	err = t.intgrSrv.Stop()
	if err != nil {
		log.Printf("intgrSrv stop err: %s", err.Error())
	}

	err = t.agentSrv.Stop()
	if err != nil {
		log.Printf("agentSrv stop err: %s", err.Error())
	}

	return nil

}

func (t *application) Start() error {

	t.exit = make(chan struct{})

	logg := microInfrastructure.NewMyUglyLogger("")
	defer logg.Close()

	fs := flag.NewFlagSet("configuration", 1)

	crypt := fs.Bool("wtf", false, "crypt")

	//app := microInfrastructure.NewWebApp(logg, version, fs)

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("parse args err: %w", err)
	}

	settRep := microInfrastructure.NewSettingsFile("settings.bin", *crypt)

	dbHand, err := infrastructure.NewSqliteHandler(t.appname + ".sqlite")
	if err != nil {
		return fmt.Errorf("NewSqliteHandler err: %w", err)
	}
	t.dbhand = dbHand

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	storer, err := interfaces.NewDataRepository(dbHand, logg)
	if err != nil {
		return fmt.Errorf("NewDataRepository err: %w", err)
	}

	agentZip, err := infrastructure.GetAgentZipsFromDir("agentzip")
	if err != nil {
		return fmt.Errorf("GetAgentZipsFromDir err: %w", err)
	}

	source, err := usecases.NewKSBScannerSource(&settRep, microDomain.About{
		Name:    t.appname,
		Version: t.appversion,
	}, logg, storer, httpClient, agentZip)
	if err != nil {
		return fmt.Errorf("NewKSBScannerSource err: %w", err)
	}

	intgrSrv, err := initIntgrSrvService(3003, source, logg)
	if err != nil {
		return fmt.Errorf("initIntgrSrvService err: %w", err)
	}
	t.intgrSrv = intgrSrv

	agentSrv, err := initAgentService(3004, source, logg)
	if err != nil {
		return fmt.Errorf("initAgentService err: %w", err)
	}
	t.agentSrv = agentSrv

	intgrErr := t.intgrSrv.Start()
	agentErr := t.agentSrv.Start()

	go func() {
		for err := range intgrErr {
			logg.Err(err, "intgrSrv err")
		}
	}()

	go func() {
		for err := range agentErr {
			logg.Err(err, "agentSrv err")
		}
	}()

	return nil

}

func initIntgrSrvService(port int, source microDomain.Source, logg microInterfaces.Logger) (*webapp.WebApp, error) {
	intgrSrv := microInterfaces.MicroServiceV1{
		JWTParams: microInterfaces.JWTParams{
			Username:        "admin",
			Password:        "123456",
			LifetimeMinutes: 5,
			Secret:          []byte(`Кикабидзе`),
		},
		Logger: logg,
		Interactor: microUsecases.NewInteractor(
			source,
			microUsecases.About{
				Name:    source.About().Name,
				Version: source.About().Version,
			}),
	}

	optionsIntgrSrv := []webapp.Option{
		webapp.Port(port),
		webapp.Logging(logg),
		//webapp.TLS(),
		webapp.Router(func(r *chi.Mux) {
			r.Route("/api/v1", func(r chi.Router) {
				r.Use(middleware.Timeout(1 * time.Second))
				r.Get("/about", intgrSrv.About)
				r.Get("/token", intgrSrv.GetToken)
				r.Group(func(r chi.Router) {
					r.Use(middleware.Timeout(60 * time.Second))
					//r.Use(intgrSrv.JWTMiddleware)
					r.Post("/assets", intgrSrv.StartAssets)
					r.Get("/assets", intgrSrv.GetAssets)
					r.Get("/incidents", intgrSrv.Incidents)
					r.Get("/settings", intgrSrv.Settings)
					r.Post("/settings", intgrSrv.SetSettings)
				})
			})
		}),
	}

	return webapp.New(optionsIntgrSrv...)

}

func initAgentService(port int, source *usecases.KSBScannerSource, logg microInterfaces.Logger) (*webapp.WebApp, error) {

	agentSrv, err := interfaces.NewAgentService(logg, source)
	if err != nil {
		return nil, fmt.Errorf("NewAgentService err: %w", err)
	}

	optionsAgentSrv := []webapp.Option{
		webapp.Port(port),
		webapp.Logging(logg),
		//webapp.TLS(),
		webapp.Router(func(r *chi.Mux) {
			r.Route("/agent/api/v1", func(r chi.Router) {
				r.Use(middleware.Timeout(1 * time.Second))
				//r.Get("/about", agentSrv.About)
				//r.Get("/token", agentSrv.GetToken)
				r.Group(func(r chi.Router) {
					r.Use(middleware.Timeout(60 * time.Second))
					//r.Use(agentSrv.JWTMiddleware)
					r.Post("/data", agentSrv.AddData)
					r.Post("/ping", agentSrv.Ping)

				})
			})
		}),
	}

	return webapp.New(optionsAgentSrv...)
	// if err != nil {
	// 	logg.Err(err, "webapp.New")
	// 	return
	// }

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)

	// errChan := web.Start()

	// select {
	// case err := <-errChan:
	// 	logg.Err(err, "webAgent.Start")
	// case sig := <-quit:
	// 	logg.Inf("Receive signal: %s", sig.String())
	// }

	// web.Stop()
}
