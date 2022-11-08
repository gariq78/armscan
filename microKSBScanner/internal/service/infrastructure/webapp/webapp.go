package webapp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
)

type WebApp struct {
	config config
	logErr func(error, string, ...interface{})
	logInf func(string, ...interface{})

	mux     *chi.Mux
	errChan chan error

	m      sync.Mutex
	server *http.Server
}

type config struct {
	port int
	cert *tls.Certificate
}

func New(options ...Option) (*WebApp, error) {
	wa := &WebApp{
		logErr: func(error, string, ...interface{}) {},
		logInf: func(string, ...interface{}) {},

		mux:     chi.NewRouter(),
		errChan: make(chan error, 1),
	}

	for _, option := range options {
		err := option(wa)
		if err != nil {
			return nil, fmt.Errorf("option: %w", err)
		}
	}

	return wa, nil
}

func (w *WebApp) Start() chan error {
	if w.server != nil {
		w.errChan <- ErrAlreadyStarted
		return w.errChan
	}

	w.m.Lock()
	w.server = &http.Server{
		Addr:           ":" + strconv.Itoa(w.config.port),
		Handler:        w.mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	w.m.Unlock()

	go func() {
		if w.config.cert != nil {
			w.logInf("start listening %d https port", w.config.port)

			w.server.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{*w.config.cert},
			}

			if err := w.server.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
				w.errChan <- fmt.Errorf("ListenAndServeTLS error %w", err)
			}
		} else {
			w.logInf("start listening %d http port", w.config.port)

			if err := w.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				w.errChan <- fmt.Errorf("ListenAndServe error %w", err)
			}
		}
	}()

	return w.errChan
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, os.Interrupt)

	// select {
	// case err := <-errChan:
	// 	w.logErr(err, "")
	// case sig := <-quit:
	// 	w.logInf("Receive signal: %s", sig.String())
	// }

	// ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	// defer shutdown()

	// return w.server.Shutdown(ctx)
}

func (w *WebApp) Stop() error {
	if w.server == nil {
		return nil
	}

	defer func() {
		close(w.errChan)
		w.server = nil
	}()

	w.logInf("stop listening %d port", w.config.port)

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	err := w.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server.Shutdown err: %w", err)
	}

	return nil
}

// func (w *WebApp) logErr(err error, msg string, args ...interface{}) {
// 	if w.logg != nil {
// 		w.logg.Err(err, msg, args...)
// 	}
// }

// func (w *WebApp) logInf(msg string, args ...interface{}) {
// 	if w.logg != nil {
// 		w.logg.Inf(msg, args...)
// 	}
// }
