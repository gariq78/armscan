package pcAntiVirusInfo

import (
	"errors"
	"fmt"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/additionInfo/kis"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/pcAntiVirusInfo/aviStructs"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/pcAntiVirusInfo/drewb"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/interfaces/pcAntiVirusInfo/nod32"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
	"os"

	"golang.org/x/sys/windows/registry"
	"log"
)

const (
	uninstallRegistryPath  = `Software\Microsoft\Windows\CurrentVersion\Uninstall` //
	uninstallRegistryPath2 = `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`
	antivirusRegistryPath  = `SOFTWARE\Microsoft\Security Center\Provider\Av`
)

var statuses = [4]string{"OFF", "ON", "SNOOZED", "EXPIRED"}

/*var avsNames =[3]string{"Kaspersky Internet Security","ESET Security","Dr.Web Security Space"}*/
type AvInfo struct {
}

func NewAVInfo() *AvInfo {
	return &AvInfo{}
}
func (i *AvInfo) GetInfo() (avs []aviStructs.AntivirusInfo, err error) {

	antivirusKey, err := registry.OpenKey(registry.LOCAL_MACHINE, antivirusRegistryPath, registry.READ)
	if err != nil {
		log.Println("GetInfo", err)
	}
	defer antivirusKey.Close()

	info, err := antivirusKey.Stat()

	if err != nil {
		// win7 нет такого key
		return nil, fmt.Errorf("uninstallKey.Stat err: %w", err)
	}
	subkeynames, err := antivirusKey.ReadSubKeyNames(int(info.SubKeyCount))
	if err != nil {
		return nil, fmt.Errorf("uninstallKey.ReadSubKeyNames err: %w", err)
	}
	for _, name := range subkeynames {
		sub, err := registry.OpenKey(registry.LOCAL_MACHINE, antivirusRegistryPath+`\`+name, registry.READ)
		if err != nil {
			return nil, fmt.Errorf("Open [%s] registry key err: %w", name, err)
		}
		var av aviStructs.AntivirusInfo
		productexe, _, _ := sub.GetStringValue("PRODUCTEXE")
		av.Name, _, _ = sub.GetStringValue("DisplayName")
		state, _, _ := sub.GetIntegerValue("state")
		// проверим наличие файла
		if i.isHasExeFile(productexe) {
			av.State = i.getState(state)
			av.SignatureStatus = i.getSignatureStatus(state)
			days, err := i.getExpiration(av)
			if err != nil {

			}
			av.Expiration = days
			av.Version, err = i.getVersionInfo(av)
			if err != nil || av.Version == "" {
				av.Version = "undefined"
				log.Println(err)
			}
			avs = append(avs, av)

		}

	}

	return

}

func (i *AvInfo) getExpiration(av aviStructs.AntivirusInfo) (days string, err error) {
	days = "undefined"

	if av.Name == "Kaspersky Internet Security" {
		var k kis.KIS
		return k.GetKISExpiration()
	}
	if av.Name == "ESET Security" {
		var n nod32.Nod32
		return n.GetExpiration(), nil
	}
	if av.Name == "Dr.Web Security Space" {
		var dr drewb.DrWeb
		avPath, _ := dr.GetInstallPath()
		return dr.GetExpiration(avPath), nil
	}
	return
}

func (i *AvInfo) getState(u64 uint64) (state string) {
	state = "undefined"
	s := u64 >> 12 & 0x0f
	if s < 4 {
		state = statuses[s]
	}

	return
}

//актульностьт анитивирусных баз
func (i *AvInfo) getSignatureStatus(u64 uint64) (status string) {
	//	UpToDate     = 0x00
	//  OutOfDate    = 0x10
	status = "undefined"
	s := u64 & 0xff
	if s == 16 {
		status = "OutOfDate"
	}
	if s == 0 {
		status = "UpToDate"
	}
	return
}

func (i *AvInfo) isHasExeFile(fileUrl string) (isHas bool) {
	_, err := os.Stat(fileUrl)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func (i *AvInfo) getVersionInfo(av aviStructs.AntivirusInfo) (version string, err error) {
	var regPath string = uninstallRegistryPath
	if av.Name == "Kaspersky Internet Security" {
		regPath = uninstallRegistryPath2
	}
	if av.Name == "ESET Security" {
		regPath = uninstallRegistryPath
	}
	if av.Name == "Dr.Web Security Space" {
		regPath = uninstallRegistryPath
	}

	uninstallKey, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.READ)
	if err != nil {
		return "undefined", fmt.Errorf("registry.OpenKey uninstall err: %w", err)
	}
	defer uninstallKey.Close()

	info, err := uninstallKey.Stat()
	if err != nil {
		fmt.Errorf("uninstallKey.Stat err: %w", err)
	}

	subkeynames, err := uninstallKey.ReadSubKeyNames(int(info.SubKeyCount))
	if err != nil {
		return "undefined", fmt.Errorf("uninstallKey.ReadSubKeyNames err: %w", err)
	}

	for _, name := range subkeynames {
		sub, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath+`\`+name, registry.READ)
		defer sub.Close()
		if err != nil {
			return "undefined", fmt.Errorf("Open [%s] registry key err: %w", name, err)
		}

		dispName, _, _ := sub.GetStringValue("DisplayName")
		if dispName == av.Name {
			version, _, _ = sub.GetStringValue("DisplayVersion")
			break

		}

	}

	return
}

func (i *AvInfo) GetAntivirusesFromSoftWareList(sws []asset.SoftwareType) (avs []aviStructs.AntivirusInfo, err error) {

	for _, sw := range sws {

		if sw.Name == "Kaspersky Internet Security" {
			var k kis.KIS
			av := i.makeAvInfo(sw)
			expDays, err := k.GetKISExpiration()
			if err != nil {

			} else {
				av.Expiration = expDays
			}

			avs = append(avs, av)
		}

		if sw.Name == "ESET Security" {

			avs = append(avs, i.makeAvInfo(sw))
		}

		if sw.Name == "Dr.Web Security Space" {
			var dr drewb.DrWeb
			path, _ := dr.GetInstallPath()
			av := i.makeAvInfo(sw)
			av.Expiration = dr.GetExpiration(path)
			avs = append(avs, av)
		}

	}

	return
}

func (i *AvInfo) makeAvInfo(sw asset.SoftwareType) (av aviStructs.AntivirusInfo) {
	av.Name = sw.Name
	av.Version = sw.Version
	av.Expiration = "not determined"
	av.SignatureStatus = "not determined"
	av.State = "not determined"
	return
}
