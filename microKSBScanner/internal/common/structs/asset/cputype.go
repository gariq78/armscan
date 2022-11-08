package asset

type CPUType struct {
	Name string `json:"name"`
}

func (a CPUType) Changed(b CPUType) bool {
	return a.Name != b.Name
}

func cpuTypesChanged(a, b []CPUType) bool {
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
