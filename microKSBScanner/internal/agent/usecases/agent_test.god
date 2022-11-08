package agent

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/domain"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/usecases/settings"
)

type tMega struct {
	inf string
	err string

	scanData  domain.Data
	scanError error

	serviceServerResponse ServiceResponce
	serviceSendDataError  error

	serviceApplySettingsError error

	saveDataError error

	settings settings.Settings
}

func (m *tMega) Inf(string, ...interface{}) {

}

func (m *tMega) Err(error, string, ...interface{}) {

}

func (m *tMega) LoadSettings() (settings.Settings, error) {
	return m.settings, nil
}

func (m *tMega) SaveSettings(s settings.Settings) error {
	m.settings = s
	return nil
}

func (m *tMega) LoadData() (domain.Data, error) {
	return domain.Data{}, nil
}

func (m *tMega) SaveData(domain.Data) error {
	if m.saveDataError != nil {
		return m.saveDataError
	}
	return nil
}

func (m *tMega) Update([]byte) error {
	return nil
}

func (m *tMega) ApplySettings(settings.Settings) error {
	if m.serviceApplySettingsError != nil {
		return m.serviceApplySettingsError
	}
	return nil
}

func (m *tMega) SendData(sett settings.Settings, data domain.Data) (ServiceResponce, error) {
	var res = ServiceResponce{}

	if m.serviceSendDataError != nil {
		return res, m.serviceSendDataError
	}

	return ServiceResponce{}, nil
}

func (m *tMega) Scan() (domain.Data, error) {
	var res = domain.Data{}

	if m.scanData != res {
		return m.scanData, nil
	}

	if m.scanError != nil {
		return res, m.scanError
	}

	return domain.Data{}, nil
}

func (m *tMega) Ping() (ServiceResponce, error) {
	return ServiceResponce{}, nil
}

func getAgent(mega *tMega, t *testing.T) *Agent {
	mega.settings = settin{}
	age, err := New(mega, mega, mega, mega, mega)
	if err != nil {
		t.Fatalf("agent.New error: %s", err.Error())
	}

	return age
}

func TestScanAndSend_ScannerScanError(t *testing.T) {
	scanErr := errors.New("scanerror")
	mega := &tMega{
		scanError: scanErr,
	}

	agent := getAgent(mega, t)

	err := agent.scanAndSend()
	if err == nil {
		t.Fatal("want scanErr but got err==nil")
	}

	if !errors.Is(err, scanErr) {
		t.Fatalf("not a scanErr [%s]", err.Error())
	}

	// OK
}

func TestScanAndSend_SameData(t *testing.T) {
	d := domain.Data{
		Value: "A",
	}

	mega := &tMega{
		scanData: d,
	}

	agent := getAgent(mega, t)
	agent.data = d

	err := agent.scanAndSend()
	if err != nil {
		t.Errorf("want err==nil but got not nil")
	}
}

func TestScanAndSend_ServiceSendDataError(t *testing.T) {
	serviceError := errors.New("a")
	mega := &tMega{
		serviceSendDataError: serviceError,
		scanData:             domain.Data{Value: "a"},
	}

	agent := getAgent(mega, t)

	err := agent.scanAndSend()
	if err == nil {
		t.Fatal("want err but got nil")
	}

	if !errors.Is(err, serviceError) {
		t.Fatal("want serviceErr but got other error")
	}
}

func TestScanAndSend_UpdateAgentData(t *testing.T) {
	mega := &tMega{
		scanData: domain.Data{Value: "a"},
	}

	agent := getAgent(mega, t)

	err := agent.scanAndSend()
	if err != nil {
		t.Fatal("want err==nil but got err!=nil")
	}

	if agent.data != mega.scanData {
		t.Fatal("data not updated in agent")
	}
}

func TestApplySettings(t *testing.T) {
	sett := settin{
		id: 99,
	}
	mega := &tMega{
		//settings: sett,
	}

	agent := getAgent(mega, t)
	//agent.settings = mega.settings

	err := agent.applySettings(sett)
	if err != nil {
		t.Fatal("want err==nil but got err!=nil")
	}

	if agent.settings.(settin).id != sett.id {
		t.Fatalf("want updated settings but got false")
	}
}

type settin struct {
	id int
}

func (s settin) GetMonitoringDuration() time.Duration {
	return 1 * time.Second
}
func (s settin) GetServerAddress() url.URL {
	return url.URL{}
}
func (s settin) RotateServerAddresses() {

}
func (s settin) NeedUpdate(n settin) bool {
	return s.id != n.id
}
