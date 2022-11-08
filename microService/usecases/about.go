package usecases

import "ksb-dev.keysystems.local/intgrsrv/microService/domain"

type About struct {
	Name    string
	Version string
	Source  domain.About
}
