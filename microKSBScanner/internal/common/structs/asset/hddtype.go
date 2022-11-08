package asset

type HDDType struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Size         string `json:"size"`
}

func (a HDDType) Changed(b HDDType) bool {
	return a.Name != b.Name ||
		a.Manufacturer != b.Manufacturer ||
		a.Size != b.Size
}

func hddTypeChanged(a, b []HDDType) bool {
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
