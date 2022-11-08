package main

import (
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/domain"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

type settRepo struct {
	settings usecases.Settings
}

var _ usecases.SettingsRepository = (*settRepo)(nil)

func (s settRepo) Load(sett usecases.Settings) error {
	sett = s.settings
	return nil
}

func (s settRepo) Save(sett usecases.Settings) error {
	s.settings = sett
	return nil
}

type dataRepo struct {
	data asset.Asset
}

var _ usecases.DataRepository = (*dataRepo)(nil)

func (d dataRepo) LoadData() (asset.Asset, error) {
	return d.data, nil
}

func (d dataRepo) SaveData(data asset.Asset) error {
	d.data = data
	return nil
}

type sensRepo struct {
	data domain.SensitiveData
}

var _ usecases.SensitiveDataRepository = (*sensRepo)(nil)

func (s sensRepo) Load(d *domain.SensitiveData) error {
	*d = s.data
	return nil
}

func (s sensRepo) Save(d domain.SensitiveData) error {
	s.data = d
	return nil
}
