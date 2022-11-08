package domain

import "time"

type Incident struct {
	Id              string    `json:"id"`
	CreateTime      time.Time `json:"createTime"`
	DetectTime      time.Time `json:"detectTime"`
	ModifyTime      time.Time `json:"modifyTime"`
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	Type            string    `json:"type"`
	Assigned        string    `json:"assigned"`
	Severity        int       `json:"severity"`
	Description     string    `json:"description"`
	Recommendations []string  `json:"recommendations"`
	Tagrets         Targets   `json:"targets"`
	Attackers       Targets   `json:"attackers"`
	Groups          []string  `json:"groups"`
}

type Targets struct {
	Assets []AssetsShort `json:"assets"`
}

type AssetsShort struct {
	HostName   string `json:"hostName"`
	IPAddress  string `json:"ipAddress"`
	MACAddress string `json:"macAddress"`
	Name       string `json:"name"`
}
