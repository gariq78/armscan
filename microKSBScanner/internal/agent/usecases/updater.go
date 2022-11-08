package usecases

import "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"

type Updater interface {
	Update(structs.AgentArchive) error
}
