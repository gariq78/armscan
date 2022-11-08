package repository

import (
	"fmt"

	agent "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

type Data struct {
	ls LoadSaver
}

var _ agent.DataRepository = Data{}

func NewData(ls LoadSaver) (Data, error) {
	rv := Data{
		ls: ls,
	}

	if ls == nil {
		return rv, fmt.Errorf("param ls is nil")
	}

	return rv, nil
}

func (d Data) LoadData() (asset.Asset, error) {
	var rv asset.Asset

	err := d.ls.Load(&rv)
	if err != nil {
		return rv, fmt.Errorf("ls.Load err: %w", err)
	}

	return rv, nil
}

func (d Data) SaveData(n asset.Asset) error {
	return d.ls.Save(n)
}
