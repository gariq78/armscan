package webapp

import (
	"crypto/tls"
	"fmt"

	"github.com/go-chi/chi"
)

type Option func(*WebApp) error

func Port(p int) Option {
	return func(wa *WebApp) error {
		if p < 1 {
			return invalidArgumentErr("Port Option p")
		}

		wa.config.port = p

		return nil
	}
}

func TLS(certPem, keyPem []byte) Option {
	return func(wa *WebApp) error {
		if certPem == nil {
			return invalidArgumentErr("TLS Option certPem")
		}

		if keyPem == nil {
			return invalidArgumentErr("TLS Option keyPem")
		}

		cert, err := tls.X509KeyPair(certPem, keyPem)
		if err != nil {
			return fmt.Errorf("TLS Option X509KeyPair: %w", err)
		}

		wa.config.cert = &cert

		return nil
	}
}

func Logging(l Logger) Option {
	return func(wa *WebApp) error {
		wa.logErr = l.Err
		wa.logInf = l.Inf

		return nil
	}
}

func Router(updateMux func(*chi.Mux)) Option {
	return func(wa *WebApp) error {
		updateMux(wa.mux)

		return nil
	}
}
