package domain

import "time"

// Source описывает источник данных
type Source interface {
	About() About
	Settings() interface{}
	SetSettings(settings interface{}) error

	Assets() (assets []Asset, err error)

	// Incidents возвращает инциденты начиная с указанного штампа и штамп времени возвращаемых инцидентов
	Incidents(notBefore time.Time) (incidents []Incident, newStamp time.Time, err error)
}
