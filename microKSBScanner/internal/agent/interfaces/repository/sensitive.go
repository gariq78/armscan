package repository

import (
	"fmt"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/domain"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
)

type SensitiveData struct {
	ls LoadSaver
}

var _ usecases.SensitiveDataRepository = SensitiveData{}

func NewSensitiveData(ls LoadSaver) (SensitiveData, error) {
	rv := SensitiveData{
		ls: ls,
	}

	if ls == nil {
		return rv, fmt.Errorf("param ls is nil")
	}

	return rv, nil
}

func (s SensitiveData) Load(to *domain.SensitiveData) error {
	return s.ls.Load(to)
}

func (s SensitiveData) Save(from domain.SensitiveData) error {
	return s.ls.Save(from)
}
