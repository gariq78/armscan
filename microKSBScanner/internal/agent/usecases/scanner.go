package usecases

import "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"

type Scanner interface {
	Scan() (asset.Asset, error)
	ID() (ID, error)
}
