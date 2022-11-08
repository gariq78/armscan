package repository

type LoadSaver interface {
	Load(interface{}) error
	Save(interface{}) error
}
