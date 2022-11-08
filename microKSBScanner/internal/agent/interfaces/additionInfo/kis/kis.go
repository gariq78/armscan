package kis

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
)

type Info struct {
	DaysTillExpiration string `json:"days_till_expiration"`
}

const (
	CROSS_PRODUCT_STORE_REG_KEY = `SOFTWARE\WOW6432Node\KasperskyLab\AVP21.3\Data\CrossProductStore`
	STR_DAYS_TILL_EXPIRATION    = `act-lic-daystillexpiration`
)

func GetInfo() (info Info, err error) {
	info, err = getInfo()
	return
}

func getInfo() (info Info, err error) {

	antivirusKey, err := registry.OpenKey(registry.LOCAL_MACHINE, CROSS_PRODUCT_STORE_REG_KEY, registry.READ)
	defer antivirusKey.Close()
	if err != nil {

		return
	}
	val, _, _ := antivirusKey.GetStringValue("SubscriptionUrl")

	n := strings.Index(val, "act-lic")
	val = val[n:]
	fmt.Println(info)

	ss := strings.Split(val, "&")

	for _, s := range ss {
		sk := strings.Split(s, "=")
		if sk[0] == STR_DAYS_TILL_EXPIRATION {
			info.DaysTillExpiration = sk[1]
		}

	}
	return
}
