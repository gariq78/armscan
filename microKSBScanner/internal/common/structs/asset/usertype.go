package asset

type UserType struct {
	Login       string `json:"login"`
	FIO         string `json:"fio"`
	Description string `json:"description"`
	LoginDate   string `json:"loginDate"`
	Domain      string `json:"domain"`
}

func (a UserType) Changed(b UserType) bool {
	return a.Login != b.Login ||
		a.FIO != b.FIO ||
		a.Description != b.Description ||
		a.LoginDate != b.LoginDate ||
		a.Domain != b.Domain
}

func userTypeChanged(a, b []UserType) bool {
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
