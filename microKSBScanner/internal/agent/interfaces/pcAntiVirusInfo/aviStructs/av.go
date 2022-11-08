package aviStructs

type AntivirusInfo struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Expiration      string `json:"expiration"`
	State           string `json:"state"`
	SignatureStatus string `json:"signature_status"`
}
