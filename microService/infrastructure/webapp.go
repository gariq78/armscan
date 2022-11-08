package infrastructure

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"ksb-dev.keysystems.local/intgrsrv/microService/domain"
	"ksb-dev.keysystems.local/intgrsrv/microService/interfaces"
	"ksb-dev.keysystems.local/intgrsrv/microService/usecases"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type WebApp struct {
	config  config
	logg    Logger
	version string
	args    args
}

type Logger interface {
	Err(error, string, ...interface{})
	Inf(string, ...interface{})
	InitSentry(string)
}

type config struct {
	AppName   string
	Port      int
	CertPem   []byte
	KeyPem    []byte
	SentryDsn string
}

type args struct {
	port      string
	cert      string
	key       string
	sentryDsn string
}

func NewWebApp(logg Logger, version string, fs *flag.FlagSet) *WebApp {
	wa := &WebApp{
		logg: logg,
		config: config{
			Port: 3000,
		},
		version: version,
		args:    args{},
	}

	fs.StringVar(&wa.args.port, "p", "", "listen port (env PORT)")
	fs.StringVar(&wa.args.cert, "cert", "", "certPem file name for https (env CERT)")
	fs.StringVar(&wa.args.key, "key", "", "keyPem file name for https (env KEY)")
	fs.StringVar(&wa.args.sentryDsn, "sentry", "", "sentry dsn (env SENTRY_DSN)")

	return wa
}

func (w *WebApp) Start(source domain.Source, jwtParams interfaces.JWTParams) {
	w.processArgs()

	if w.config.SentryDsn != "" {
		//w.logg = newLogger(w.config.SentryDsn)
		w.logg.InitSentry(w.config.SentryDsn)
	}

	srcAbout := source.About()

	ms := interfaces.MicroServiceV1{
		JWTParams: jwtParams,
		Logger:    w.logg,
		Interactor: usecases.NewInteractor(
			source,
			usecases.About{
				Name:    "micro" + srcAbout.Name,
				Version: w.version,
			}),
	}

	err := w.startWeb(ms)
	if err != nil {
		w.logg.Err(err, "shutdown server error")
	}
}

func (w *WebApp) startWeb(ms interfaces.MicroServiceV1) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Timeout(1 * time.Second))
		r.Get("/about", ms.About)
		r.Get("/token", ms.GetToken)
		//r.Get("/status", ms.Status)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Timeout(60 * time.Second))
			r.Use(ms.JWTMiddleware)
			r.Post("/assets", ms.StartAssets)
			r.Get("/assets", ms.GetAssets)
			r.Get("/incidents", ms.Incidents)
			r.Get("/settings", ms.Settings)
			r.Post("/settings", ms.SetSettings)
		})
	})

	// HTTP Server
	httpServer := &http.Server{
		Addr:           ":" + strconv.Itoa(w.config.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if len(w.config.CertPem) > 0 && len(w.config.KeyPem) > 0 {
		cert, err := tls.X509KeyPair(w.config.CertPem, w.config.KeyPem)
		if err != nil {
			return fmt.Errorf("X509KeyPair: %w", err)
		}

		httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	errChan := make(chan error, 1)
	go func() {
		if httpServer.TLSConfig != nil {
			w.logg.Inf("start listening %d https port", w.config.Port)
			if err := httpServer.ListenAndServeTLS("", ""); err != nil {
				errChan <- fmt.Errorf("ListenAndServeTLS error %w", err)
			}
		} else {
			w.logg.Inf("start listening %d http port", w.config.Port)
			if err := httpServer.ListenAndServe(); err != nil {
				errChan <- fmt.Errorf("ListenAndServe error %w", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	select {
	case err := <-errChan:
		w.logg.Err(err, "")
	case sig := <-quit:
		w.logg.Inf("Receive signal: %s", sig.String())
	}

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return httpServer.Shutdown(ctx)
}

func (w *WebApp) processArgs() {
	if w.args.port == "" {
		w.args.port, _ = os.LookupEnv("PORT")
	}
	if w.args.port != "" {
		port, err := strconv.Atoi(w.args.port)
		if err != nil {
			w.logg.Err(err, "error converting PORT: %s", w.args.port)
		} else {
			w.config.Port = port
		}
	}

	if w.args.cert == "" {
		w.args.cert, _ = os.LookupEnv("CERT")
	}
	if w.args.cert != "" {
		b, err := ioutil.ReadFile(w.args.cert)
		if err != nil {
			w.logg.Err(err, "read cert file [%s]", w.args.cert)
		} else {
			w.config.CertPem = b
		}
	}

	if w.args.key == "" {
		w.args.key, _ = os.LookupEnv("KEY")
	}
	if w.args.key != "" {
		b, err := ioutil.ReadFile(w.args.key)
		if err != nil {
			w.logg.Err(err, "read key file [%s]", w.args.key)
		} else {
			w.config.KeyPem = b
		}
	}

	if w.args.sentryDsn == "" {
		w.args.sentryDsn, _ = os.LookupEnv("SENTRY_DSN")
	}
	if w.args.sentryDsn != "" {
		w.config.SentryDsn = w.args.sentryDsn
	}
}
