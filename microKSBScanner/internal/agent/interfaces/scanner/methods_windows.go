package scanner

import (
	"fmt"

	"golang.org/x/sys/windows/registry"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

const uninstallRegistryPath = `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`
const uninstallRegistryPath2 = `Software\Microsoft\Windows\CurrentVersion\Uninstall` //

func getSoftwareList() (softwareList []asset.SoftwareType, err error) {
	softwareList, err = getSoftwareListFromRegPath(uninstallRegistryPath)
	softwareList2, err := getSoftwareListFromRegPath(uninstallRegistryPath2)
	softwareList = append(softwareList, softwareList2...)
	softwareList = removeDuplicates(softwareList)
	return
}

func getSoftwareListFromRegPath(regPath string) ([]asset.SoftwareType, error) {

	uninstallKey, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("registry.OpenKey uninstall err: %w", err)
	}
	defer uninstallKey.Close()

	info, err := uninstallKey.Stat()
	if err != nil {
		return nil, fmt.Errorf("uninstallKey.Stat err: %w", err)
	}

	subkeynames, err := uninstallKey.ReadSubKeyNames(int(info.SubKeyCount))
	if err != nil {
		return nil, fmt.Errorf("uninstallKey.ReadSubKeyNames err: %w", err)
	}

	rv := make([]asset.SoftwareType, 0, info.SubKeyCount)

	for _, name := range subkeynames {
		sub, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath+`\`+name, registry.READ)
		if err != nil {
			return nil, fmt.Errorf("Open [%s] registry key err: %w", name, err)
		}

		dispName, _, _ := sub.GetStringValue("DisplayName")
		dispVers, _, _ := sub.GetStringValue("DisplayVersion")
		publisher, _, _ := sub.GetStringValue("Publisher")
		comments, _, _ := sub.GetStringValue("Comments")

		if dispName != "" { // в этой ветке есть инфа о обновлениях безопасности, возможно в будущем нужно будет вытащить
			rv = append(rv, asset.SoftwareType{
				Name:         dispName,
				Manufacturer: publisher,
				Version:      dispVers,
				Description:  comments,
			})
		}

		sub.Close()
	}

	return rv, nil
}

func contains(s []asset.SoftwareType, e asset.SoftwareType) bool {
	for _, a := range s {
		if a.Name == e.Name && a.Version == e.Version && a.Manufacturer == e.Manufacturer {
			return true
		}
	}
	return false
}

func removeDuplicates(strList []asset.SoftwareType) []asset.SoftwareType {
	list := []asset.SoftwareType{}
	for _, item := range strList {

		if contains(list, item) == false {
			list = append(list, item)
		}
	}
	return list
}

// old.
/*func getSoftwareList() ([]asset.SoftwareType, error) {

	uninstallKey, err := registry.OpenKey(registry.LOCAL_MACHINE, uninstallRegistryPath, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("registry.OpenKey uninstall err: %w", err)
	}
	defer uninstallKey.Close()

	info, err := uninstallKey.Stat()
	if err != nil {
		return nil, fmt.Errorf("uninstallKey.Stat err: %w", err)
	}

	subkeynames, err := uninstallKey.ReadSubKeyNames(int(info.SubKeyCount))
	if err != nil {
		return nil, fmt.Errorf("uninstallKey.ReadSubKeyNames err: %w", err)
	}

	rv := make([]asset.SoftwareType, 0, info.SubKeyCount)

	for _, name := range subkeynames {
		sub, err := registry.OpenKey(registry.LOCAL_MACHINE, uninstallRegistryPath+`\`+name, registry.READ)
		if err != nil {
			return nil, fmt.Errorf("Open [%s] registry key err: %w", name, err)
		}

		dispName, _, _ := sub.GetStringValue("DisplayName")
		dispVers, _, _ := sub.GetStringValue("DisplayVersion")
		publisher, _, _ := sub.GetStringValue("Publisher")
		comments, _, _ := sub.GetStringValue("Comments")

		if dispName != "" { // в этой ветке есть инфа о обновлениях безопасности, возможно в будущем нужно будет вытащить
			rv = append(rv, asset.SoftwareType{
				Name:         dispName,
				Manufacturer: publisher,
				Version:      dispVers,
				Description:  comments,
			})
		}

		sub.Close()
	}

	return rv, nil
}*/
