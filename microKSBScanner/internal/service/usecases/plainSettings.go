package usecases

import (
	"encoding/json"
	"strconv"
)

type plainSettings struct {
	AgentServiceAddress           string
	AgentAdditionalServiceAddress string
	AgentPingPeriodHours          int
	AgentMonitoringHours          int
}

func (s plainSettings) MarshalJSON() ([]byte, error) {
	type Alias plainSettings
	return json.Marshal(&struct {
		*Alias
		AgentPingPeriodHours string
		AgentMonitoringHours string
	}{
		Alias:                (*Alias)(&s),
		AgentPingPeriodHours: strconv.Itoa(s.AgentPingPeriodHours),
		AgentMonitoringHours: strconv.Itoa(s.AgentMonitoringHours),
	})
}

func (s *plainSettings) UnmarshalJSON(data []byte) error {
	type Alias plainSettings
	temp := &struct {
		*Alias
		AgentPingPeriodHours string
		AgentMonitoringHours string
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.AgentPingPeriodHours != "" {
		period, err := strconv.Atoi(temp.AgentPingPeriodHours)
		if err != nil {
			return err
		}

		s.AgentPingPeriodHours = period
	}

	if temp.AgentMonitoringHours != "" {
		period, err := strconv.Atoi(temp.AgentMonitoringHours)
		if err != nil {
			return err
		}

		s.AgentMonitoringHours = period
	}

	return nil
}
