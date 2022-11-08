package structs

import (
	"time"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

type AddDataRequset struct {
	Version int
	ID      string
	Request AddDataPacket
}

type AddDataPacket struct {
	PingRequestData
	ScanTime time.Time
	Asset    asset.Asset
}
