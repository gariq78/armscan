package interfaces

type intgrService struct {
}

func NewIntgrSerivce() (*intgrService, error) {
	rv := &intgrService{}

	return rv, nil
}
