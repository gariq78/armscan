package additionInfo

import (
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/additionInfo/kes"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/additionInfo/kis"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/additionInfo/vipNet"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
	"log"
	"strings"
)

var softWareNames = []string{
	"ViPNet Client",
	"Kaspersky Endpoint Security для Windows",
	"Kaspersky Internet Security",
}

func AddInfo(softWare []asset.SoftwareType) (modifySW []asset.SoftwareType) {
	for _, sw := range softWare {
		sw = addInfoToSoftWare(sw)
		if sw.AdditionInfo == "" {
			sw.AdditionInfo = "None"
		}

		modifySW = append(modifySW, sw)

	}

	return

}

func addInfoToSoftWare(softWare asset.SoftwareType) asset.SoftwareType {
	if strings.ToUpper(softWare.Name) == strings.ToUpper("Kaspersky Internet Security") {

		info, err := kis.GetInfo()
		if err != nil {
		}
		softWare.AdditionInfo = info

	}
	if strings.ToUpper(softWare.Name) == strings.ToUpper("Kaspersky Endpoint Security для Windows") {

		info, err := kes.GetInfo()
		if err != nil {
			log.Println(err)
		}

		softWare.AdditionInfo = info

	}

	if strings.ToUpper(softWare.Name) == strings.ToUpper("ViPNet Client") {
		info, err := vipNet.GetInfo()

		if err != nil {

		}
		softWare.AdditionInfo = info

	}
	return softWare
}
