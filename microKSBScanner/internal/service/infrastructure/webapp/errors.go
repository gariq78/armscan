package webapp

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyStarted = errors.New("server already started")
	ErrArgument       = errors.New("invalid argument")
)

func invalidArgumentErr(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrArgument)
}
