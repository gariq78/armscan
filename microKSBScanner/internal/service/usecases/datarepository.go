package usecases

import (
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

type DataRepository interface {
	AddData(m structs.AddDataPacket) error
	Datas() ([]asset.Asset, error)
}
