package structs

import (
	"fmt"
	"net/url"
	"time"
)

// 	GetPingDuration() time.Duration
// 	GetMonitoringDuration() time.Duration
// 	GetServerAddress() url.URL
// 	RotateServerAddresses()
// 	NeedUpdate(interface{}) bool
// 	json.Marshaler
// 	json.Unmarshaler

type AgentSettings struct {
	ServiceAddress           string
	AdditionalServiceAddress string
	PingPeriodHours          int
	MonitoringHours          int
}

func (s AgentSettings) GetPingDuration() time.Duration {
	//fmt.Printf("GetPingingDuration: %v\n", time.Duration(s.pingingEveryHours)*time.Minute)
	return time.Duration(s.PingPeriodHours) * time.Minute // TODO: time.Hour
}

func (s AgentSettings) GetMonitoringDuration() time.Duration {
	//fmt.Printf("GetMonitoringDuration: %v\n", time.Duration(s.monitoringEveryHours)*time.Minute)
	return time.Duration(s.MonitoringHours) * time.Minute // TODO: time.Hour
}

func (s AgentSettings) GetServerAddress() url.URL {
	rv, err := url.Parse(s.ServiceAddress)
	if err != nil {
		panic(fmt.Sprintf("cannot parse ServiceAddress %s", s.ServiceAddress))
	}
	return *rv
}

func (s *AgentSettings) RotateServerAddresses() {
	tmp := s.ServiceAddress
	s.ServiceAddress = s.AdditionalServiceAddress
	s.AdditionalServiceAddress = tmp
}

func (s AgentSettings) NeedUpdate(n interface{}) bool {
	if ns, ok := n.(AgentSettings); ok {
		if s.PingPeriodHours != ns.PingPeriodHours {
			//fmt.Printf("NeedUpdate: pingingEveryHours %d!=%d\n", s.pingingEveryHours, ns.pingingEveryHours)
			return true
		}

		if s.MonitoringHours != ns.MonitoringHours {
			//fmt.Printf("NeedUpdate: monitoringEveryHours %d!=%d\n", s.monitoringEveryHours, ns.monitoringEveryHours)
			return true
		}

		if s.ServiceAddress != ns.ServiceAddress {
			return true
		}

		if s.AdditionalServiceAddress != ns.AdditionalServiceAddress {
			return true
		}

	} else {
		//fmt.Printf("NeedUpdate: n is not a AgentSettings, but is %T\n", n)
	}

	return false
}
