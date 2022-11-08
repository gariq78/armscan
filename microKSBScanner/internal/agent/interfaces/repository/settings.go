package repository

import (
	"fmt"

	agent "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
)

type Settings struct {
	ls LoadSaver
}

var _ agent.SettingsRepository = Settings{}

func NewSettings(ls LoadSaver) (Settings, error) {
	rv := Settings{
		ls: ls,
	}

	if ls == nil {
		return rv, fmt.Errorf("param ls is nil")
	}

	return rv, nil
}

func (s Settings) Load(to *structs.AgentSettings) error {
	return s.ls.Load(to)
}

func (s Settings) Save(from structs.AgentSettings) error {
	return s.ls.Save(from)
}
