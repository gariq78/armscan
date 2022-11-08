package asset

type Asset struct {
	ClientID      string         `json:"clientid"`
	HostID        string         `json:"hostid"`
	HostName      string         `json:"hostName"`
	Domain        string         `json:"domain"`
	IPAddress     string         `json:"ipAddress"`
	MACAddress    string         `json:"macAddress"`
	Name          string         `json:"name"`
	OSName        string         `json:"osName"`
	OSVersion     string         `json:"osVersion"`
	CPU           []CPUType      `json:"cpu"`
	SystemMemory  string         `json:"systemMemory"`
	Video         string         `json:"video"`
	HDD           []HDDType      `json:"hdd"`
	FactoryNumber string         `json:"factoryNumber"`
	Motherboard   string         `json:"motherboard"`
	Monitor       string         `json:"monitor"`
	OD            string         `json:"od"`
	Software      []SoftwareType `json:"software"`
	Users         []UserType     `json:"users"`
}

func (d Asset) Changed(n Asset) bool {
	return d.ClientID != n.ClientID ||
		d.HostID != n.HostID ||
		d.HostName != n.HostName ||
		d.Domain != n.Domain ||
		d.IPAddress != n.IPAddress ||
		d.MACAddress != n.MACAddress ||
		d.Name != n.Name ||
		d.OSName != n.OSName ||
		d.OSVersion != n.OSVersion ||
		d.SystemMemory != n.SystemMemory ||
		d.Video != n.Video ||
		d.FactoryNumber != n.FactoryNumber ||
		d.Motherboard != n.Motherboard ||
		d.Monitor != n.Monitor ||
		d.OD != n.OD ||
		cpuTypesChanged(d.CPU, n.CPU) ||
		hddTypeChanged(d.HDD, n.HDD) ||
		softwareTypeChanged(d.Software, n.Software) ||
		userTypeChanged(d.Users, n.Users)
}
