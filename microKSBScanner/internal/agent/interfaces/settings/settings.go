package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

/*
type Settings interface {
	GetMonitoringDuration() time.Duration
	GetServerAddress() url.URL
	RotateServerAddresses()
	NeedUpdate(Settings) bool
}
*/
//var _ agent.Settings = &Settings{}

type Settings struct {
	valid                bool
	pingingEveryHours    int
	monitoringEveryHours int
	serverAddresses      []*url.URL
}

func New(pingingEveryHours int, monitoringEveryHours int, serverAddresses []string) (Settings, error) {
	var res Settings
	if monitoringEveryHours <= 0 {
		return res, errors.New("monitoringEveryHours cannot be less or equal zero")
	}

	res.monitoringEveryHours = monitoringEveryHours

	if pingingEveryHours <= 0 {
		return res, errors.New("pingingEveryHours cannot be less or equal zero")
	}

	res.pingingEveryHours = pingingEveryHours

	err := res.fillServerAddresses(serverAddresses)
	if err != nil {
		return res, fmt.Errorf("fillServerAddresses err: %w", err)
	}

	return res, nil
}

func (s *Settings) fillServerAddresses(serverAddresses []string) error {
	// Закоментарил так как при передаче настройки через название файла возможен только один адрес
	// if len(serverAddresses) < 2 {
	// 	return errors.New("at least two serverAddresses must be specified")
	// }

	s.serverAddresses = make([]*url.URL, 0, len(serverAddresses))

	for _, addr := range serverAddresses {
		u, err := url.ParseRequestURI(addr)
		if err != nil {
			return fmt.Errorf("parse address [%s] error: %w", addr, err)
		}

		s.serverAddresses = append(s.serverAddresses, u)
	}

	return nil
}

func (s Settings) GetPingDuration() time.Duration {
	//fmt.Printf("GetPingingDuration: %v\n", time.Duration(s.pingingEveryHours)*time.Minute)
	return time.Duration(s.pingingEveryHours) * time.Minute // TODO: time.Hour
}

func (s Settings) GetMonitoringDuration() time.Duration {
	//fmt.Printf("GetMonitoringDuration: %v\n", time.Duration(s.monitoringEveryHours)*time.Minute)
	return time.Duration(s.monitoringEveryHours) * time.Minute // TODO: time.Hour
}

func (s Settings) GetServerAddress() url.URL {
	return *s.serverAddresses[0]
}

func (s Settings) RotateServerAddresses() {
	s.serverAddresses = rotateUrls(s.serverAddresses, 1)
}

func (s Settings) NeedUpdate(n interface{}) bool {
	if ns, ok := n.(*Settings); ok {
		if s.pingingEveryHours != ns.pingingEveryHours {
			//fmt.Printf("NeedUpdate: pingingEveryHours %d!=%d\n", s.pingingEveryHours, ns.pingingEveryHours)
			return true
		}

		if s.monitoringEveryHours != ns.monitoringEveryHours {
			//fmt.Printf("NeedUpdate: monitoringEveryHours %d!=%d\n", s.monitoringEveryHours, ns.monitoringEveryHours)
			return true
		}

		if len(s.serverAddresses) != len(ns.serverAddresses) {
			//fmt.Printf("NeedUpdate: len serverAddresses %d!=%d\n", len(s.serverAddresses), len(ns.serverAddresses))
			return true
		}

		for i := range s.serverAddresses {
			if s.serverAddresses[i].String() != ns.serverAddresses[i].String() {
				//fmt.Printf("NeedUpdate: serverAddresses %s!=%s\n", s.serverAddresses[i].String(), ns.serverAddresses[i].String())
				return true
			}
		}
	} else {
		//fmt.Printf("NeedUpdate: n is not a Settings, but is %T\n", n)
	}

	return false
}

type tempType struct {
	PingingEveryHours    int
	MonitoringEveryHours int
	ServerAddresses      []string
}

func (s Settings) MarshalJSON() ([]byte, error) {
	rv := make([]string, 0, len(s.serverAddresses))
	for _, u := range s.serverAddresses {
		rv = append(rv, u.String())
	}

	var tmp = tempType{
		PingingEveryHours:    s.pingingEveryHours,
		MonitoringEveryHours: s.monitoringEveryHours,
		ServerAddresses:      rv,
	}

	return json.Marshal(tmp)
}

func (s *Settings) UnmarshalJSON(b []byte) error {
	var tmp tempType

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	s.pingingEveryHours = tmp.PingingEveryHours
	s.monitoringEveryHours = tmp.MonitoringEveryHours

	return s.fillServerAddresses(tmp.ServerAddresses)
}

func rotateUrls(urls []*url.URL, k int) []*url.URL {
	if k < 0 || len(urls) == 0 {
		return urls
	}

	r := len(urls) - k%len(urls)
	urls = append(urls[r:], urls[:r]...)

	return urls
}
