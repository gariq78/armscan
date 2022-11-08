package installer

import (
	"fmt"
	"strings"
)

var serverTemplate = `http://host:port/agent/api/v1`

// decodeFileName decodes values inventory_81fb3e6c54c44f728c9d84c864a96b7b_corporate.local_4000.exe
func decodeFileName(name string) (clientID, serverAddress string, err error) {
	arr := strings.Split(name, "_")
	if len(arr) != 4 {
		err = fmt.Errorf("bad name")
		return
	}

	clientID = arr[1]
	serverAddress = strings.Replace(serverTemplate, "host:port", arr[2]+":"+strings.Replace(arr[3], ".exe", "", 1), 1)
	return
}
