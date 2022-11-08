package usecases

import (
	"sync"
	"time"

	"ksb-dev.keysystems.local/intgrsrv/microService/domain"
)

type Interactor struct {
	mu     sync.Mutex
	busy   bool
	assets []domain.Asset
	err    error
	info   About
	source domain.Source
}

func NewInteractor(src domain.Source, about About) *Interactor {
	about.Source = src.About()

	return &Interactor{
		source: src,
		info:   about,
	}
}

func (m *Interactor) About() About {
	return m.info
}

// StartAssets инициирует запрос активов из источника данных
// возвращает true если запустили, иначе уже было запущено
func (m *Interactor) StartAssets() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.busy {
		go func() {
			assets, err := m.source.Assets()
			m.mu.Lock()
			defer m.mu.Unlock()

			if err != nil {
				m.err = err
				m.assets = nil
			} else {
				m.err = nil
				m.assets = assets
			}

			m.busy = false
		}()

		m.busy = true
		return true
	}

	return false
}

// GetAssets возращает активы источника данных
// если массив = nil значит данные еще не получены либо в процессе получения
func (m *Interactor) GetAssets() ([]domain.Asset, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.busy {
		return nil, Busy
	}
	return m.source.Assets()
	//return m.assets, m.err
}

// Incidents возвращает инциденты начиная с переданного штампа времени
func (m *Interactor) Incidents(p time.Time) ([]domain.Incident, time.Time, error) {
	return m.source.Incidents(p)
}

// Settings возвращает настройки источника
func (m *Interactor) Settings() interface{} {
	return m.source.Settings()
}

// SetSettings передает настройки источнику
func (m *Interactor) SetSettings(s interface{}) error {
	return m.source.SetSettings(s)
}
