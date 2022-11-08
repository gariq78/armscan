package infrastructure

import (
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type MyUglyLogger struct {
	// чето не очень решение
	err   func(error, string, ...interface{})
	inf   func(string, ...interface{})
	close func() error
}

func NewMyUglyLogger(sentryDsn string) *MyUglyLogger {
	return &MyUglyLogger{
		err: func(err error, msg string, args ...interface{}) {
			t := fmt.Sprintf(msg, args...)
			log.Printf("err: %s [%s]\n", t, err)
		},
		inf: func(msg string, args ...interface{}) {
			log.Printf("inf: "+msg+"\n", args...)
		},
		close: func() error {
			return nil
		},
	}
}

func (l *MyUglyLogger) Err(err error, msg string, args ...interface{}) {
	l.err(err, msg, args...)
}

func (l *MyUglyLogger) Inf(msg string, args ...interface{}) {
	l.inf(msg, args...)
}

func (l *MyUglyLogger) Close() error {
	return l.close()
}

func (l *MyUglyLogger) InitSentry(sentryDsn string) {
	if sentryDsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:   sentryDsn,
			Debug: true,
		})
		if err != nil {
			l.Err(err, "sentry.Init dsn: [%s]", sentryDsn)
			return
		}

		*l = MyUglyLogger{
			err: func(err error, msg string, args ...interface{}) {
				sentry.CaptureException(err)
			},
			inf: func(msg string, args ...interface{}) {
				sentry.CaptureMessage(fmt.Sprintf(msg, args...))
			},
			close: func() error {
				sentry.Flush(2 * time.Second)
				return nil
			},
		}
	}
}
