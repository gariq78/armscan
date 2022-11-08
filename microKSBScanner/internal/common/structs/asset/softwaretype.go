package asset

type SoftwareType struct {
	Name         string      `json:"name"`
	Manufacturer string      `json:"manufacturer"`
	Version      string      `json:"version"`
	Description  string      `json:"description"`
	AdditionInfo interface{} `json:"addition_info"`
}

func (a SoftwareType) Changed(b SoftwareType) bool {
	return a.Name != b.Name ||
		a.Manufacturer != b.Manufacturer ||
		a.Version != b.Version ||
		a.Description != b.Description ||
		a.AdditionInfo != b.AdditionInfo
}

func softwareTypeChanged(a, b []SoftwareType) bool {
	if len(a) != len(b) {
		return true
	}

	for i := range a {
		if a[i].Changed(b[i]) {
			return true
		}
	}

	return false
}
