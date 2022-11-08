package usecases

import "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"

type SettingsRepository interface {
	Load(interface{}) error
	Save(interface{}) error
}

type Settings struct {
	Agent structs.AgentSettings
}

// type agentSettings struct {
// 	ServiceAddress           string
// 	AdditionalServiceAddress string
// 	PingPeriodHours          int
// 	MonitoringHours          int
// }

// func (s agentSettings) MarshalJSON() ([]byte, error) {
// 	type Alias agentSettings
// 	return json.Marshal(&struct {
// 		*Alias
// 		PingPeriodHours string
// 		MonitoringHours string
// 	}{
// 		Alias:           (*Alias)(&s),
// 		PingPeriodHours: strconv.Itoa(s.PingPeriodHours),
// 		MonitoringHours: strconv.Itoa(s.MonitoringHours),
// 	})
// }

// func (s *agentSettings) UnmarshalJSON(data []byte) error {
// 	type Alias agentSettings
// 	temp := &struct {
// 		*Alias
// 		PingPeriodHours string
// 		MonitoringHours string
// 	}{
// 		Alias: (*Alias)(s),
// 	}

// 	if err := json.Unmarshal(data, &temp); err != nil {
// 		return err
// 	}

// 	if temp.PingPeriodHours != "" {
// 		period, err := strconv.Atoi(temp.PingPeriodHours)
// 		if err != nil {
// 			return err
// 		}

// 		s.PingPeriodHours = period
// 	}

// 	if temp.MonitoringHours != "" {
// 		period, err := strconv.Atoi(temp.MonitoringHours)
// 		if err != nil {
// 			return err
// 		}

// 		s.MonitoringHours = period
// 	}

// 	return nil
// }
