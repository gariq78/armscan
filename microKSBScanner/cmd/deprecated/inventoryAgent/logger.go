package main

import (
	"fmt"
	"os"
)

type logger struct {
	name string
	f    *os.File
}

func NewLogger(name string) (*logger, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("create log file err: %w", err)
	}

	return &logger{
		name: name,
		f:    f,
	}, nil
}

func (l *logger) Close() error {
	return l.f.Close()
}

func (l *logger) Delete() error {
	return os.Remove(l.name)
}

func (l *logger) Err(err error, msg string, args ...interface{}) {
	fmt.Fprintf(l.f, "[error] %s - %s", err.Error(), fmt.Sprintf(msg, args...))
}

func (l *logger) Inf(msg string, args ...interface{}) {
	fmt.Fprintf(l.f, "[info] %s", fmt.Sprintf(msg, args...))
}
