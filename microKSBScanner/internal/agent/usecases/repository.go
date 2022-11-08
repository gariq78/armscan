package usecases

import (
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/domain"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

type SettingsRepository interface {
	Load(*structs.AgentSettings) error
	Save(structs.AgentSettings) error
}

type DataRepository interface {
	LoadData() (asset.Asset, error)
	SaveData(asset.Asset) error
}

type SensitiveDataRepository interface {
	Load(*domain.SensitiveData) error
	Save(domain.SensitiveData) error
}
