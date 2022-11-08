package domain

type Asset struct {
	/*
		Личный идентификатор бабушки.
		Выдаётся ей при регистрации в системе для нашего сканнера.
		В МаксПатроле и Касперском видимо будет пусто.
	*/
	ClientID string `json:"clientid"`

	/*
		Хост айди генерируемое ОС
	*/
	HostID string `json:"hostid"`

	HostName      string          `json:"hostName"`
	Domain        string          `json:"domain"`
	IPAddress     string          `json:"ipAddress"`
	MACAddress    string          `json:"macAddress"`
	Name          string          `json:"name"`
	OSName        string          `json:"osName"`
	OSVersion     string          `json:"osVersion"`
	CPU           []CPUType       `json:"cpu"`
	SystemMemory  string          `json:"systemMemory"`
	Video         string          `json:"video"`
	HDD           []HDDType       `json:"hdd"`
	FactoryNumber string          `json:"factoryNumber"`
	Motherboard   string          `json:"motherboard"`
	Monitor       string          `json:"monitor"`
	OD            string          `json:"od"`
	Software      []SoftwareType  `json:"software"`
	Users         []UsersType     `json:"users"`
	Antiviruses   []AntivirusInfo `json:"antiviruses"`
}

type CPUType struct {
	Name string `json:"name"`
}

type HDDType struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Size         string `json:"size"`
}

type SoftwareType struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Version      string `json:"version"`
	Description  string `json:"description"`
}

type UsersType struct {
	Login       string `json:"login"`
	FIO         string `json:"fio"`
	Description string `json:"description"`
	LoginDate   string `json:"loginDate"`
	Domain      string `json:"domain"`
}
type AntivirusInfo struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Expiration      string `json:"expiration"`
	State           string `json:"state"`
	SignatureStatus string `json:"signature_status"`
}
