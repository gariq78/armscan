package kes

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

type KesInfoValues struct {
	REG_INSTALL_DATE_PATH                       string
	REG_DATA                                    string
	REG_TICKET                                  string
	TICKET_SUB_KEY                              string
	LAST_START_SUBKEY_NAME                      string
	LAST_SUCCESSFUL_FULL_SCAN_SUBKEY_NAME       string
	LAST_SUCCESSFUL_FULL_SCAN_START_SUBKEY_NAME string
	LAST_SUCCESSFUL_UPDATE_SUBKEY_NAME          string
	LAST_UPDATE_SUBKEY_NAME                     string
}

type Ticket struct {
	ActivationCenterId string `json:"ActivationCenterId"`
	ActivationCodeHash string `json:"ActivationCodeHash"`
	BodyHash           string `json:"BodyHash"`
	CreatedAt          string `json:"CreatedAt"`
	HardwareId         string `json:"HardwareId"`
	KPCBoundState      int    `json:"KPCBoundState"`
	KPCUnboundReason   string `json:"KPCUnboundReason"`
	LicenseId          string `json:"LicenseId"`
	LicenseState       int    `json:"LicenseState"`
	LicenseVersion     int    `json:"LicenseVersion"`
	TicketId           string `json:"TicketId"`
	TicketSequenceId   string `json:"TicketSequenceId"`
	ValidFrom          string `json:"ValidFrom"`
	ValidTo            string `json:"ValidTo"`
	Version            int    `json:"Version"`
	ExtraData          string `json:"ExtraData"`
}

type LastsData struct {
	LastStart                   int64 `json:"last_start"`
	LastSuccessfulFullScan      int64 `json:"last_successful_full_scan"`
	LastSuccessfulFullScanStart int64 `json:"last_successful_full_scan_start"`
	LastSuccessfulUpdate        int64 `json:"last_successful_update"`
	LastUpdate                  int64 `json:"last_update"`
}

type Info struct {
	InstallKESDate                  string `json:"install_kes_date"`
	LicenseExpirationDate           string `json:"license_expiration_date"`
	DaysTillLicenseExpiration       int    `json:"daysTillLicenseExpiration"`
	LastStartDate                   string `json:"last_start_date"`
	LastSuccessfulFullScanDate      string `json:"last_successful_full_scan_date"`
	LastSuccessfulFullScanStartDate string `json:"last_successful_full_scan_start_date"`
	LastSuccessfulUpdateDate        string `json:"last_successful_update_date"`
	LastUpdateDate                  string `json:"last_update_date"`
	LicenseID                       string `json:"license_id"`
}

const (
	REG_CONNECTOR_CACHE = "SOFTWARE\\WOW6432Node\\KasperskyLab\\ConnectorCache"
)

var values KesInfoValues

func GetInfo() (info Info, err error) {
	tempVal, err := getKesInfoValues()
	if err != nil {
		saveErrToFile(err)
		return
	}

	values = tempVal
	kesInstallTimeStamp, err := getInstallTimeStamp()

	infoTiket, err := getTicket()
	if err != nil {
		saveErrToFile(err)
		return
	}
	from := time.Now()
	to, _ := getTimeFromDate(infoTiket.ValidTo)

	var daysTill int
	daysTill = int(to.Sub(from).Hours() / 24)
	lastesData, err := getLastData()
	if err != nil {
		saveErrToFile(err)
		return

	}
	info.InstallKESDate = getDateFromTimeStampSec(kesInstallTimeStamp)
	info.LicenseExpirationDate = getDateFromTimeStampSec(to.Unix())
	info.DaysTillLicenseExpiration = daysTill
	info.LastStartDate = getDateFromTimeStampSec(lastesData.LastStart)
	info.LastSuccessfulFullScanDate = getDateFromTimeStampSec(lastesData.LastSuccessfulFullScan)
	info.LastSuccessfulFullScanStartDate = getDateFromTimeStampSec(lastesData.LastSuccessfulFullScanStart)
	info.LastSuccessfulUpdateDate = getDateFromTimeStampSec(lastesData.LastSuccessfulUpdate)
	info.LastUpdateDate = getDateFromTimeStampSec(lastesData.LastSuccessfulUpdate)
	info.LicenseID = strings.ToUpper(infoTiket.LicenseId)

	return

}

func getDayFromTimeStamp(stamp int64) (ts time.Time) {

	ts = time.Unix(0, stamp)

	return
}

func getInstallTimeStamp() (timeStamp int64, err error) {

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, values.REG_INSTALL_DATE_PATH, registry.READ)
	defer key.Close()
	if err != nil {
		saveErrToFile(err)
		log.Println(err)
		return 0, err
	}
	timeStampStr, _, err := key.GetStringValue("InstallDate")
	if err != nil {
		saveErrToFile(err)

		return 0, err
	}
	timeStamp, err = stringToInt64(timeStampStr)
	if err != nil {
		saveErrToFile(err)

		return 0, err

	}
	return
}

func getLastData() (lastData LastsData, err error) {

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, values.REG_DATA, registry.READ)
	defer key.Close()
	if err != nil {
		saveErrToFile(err)
		log.Println("registry.OpenKey:", err)
		return
	}
	lastStart, _, err := key.GetIntegerValue(values.LAST_START_SUBKEY_NAME)
	if err != nil {
		saveErrToFile(err)
		log.Println(err)

	}
	lastData.LastStart = int64(lastStart)

	lastSuccessfulFullScan, _, err := key.GetIntegerValue(values.LAST_SUCCESSFUL_FULL_SCAN_SUBKEY_NAME)
	if err != nil {
		saveErrToFile(err)
		log.Println(err)
	}
	lastData.LastSuccessfulFullScan = int64(lastSuccessfulFullScan)

	lastSuccessfulFullScanStart, _, err := key.GetIntegerValue(values.LAST_SUCCESSFUL_FULL_SCAN_START_SUBKEY_NAME)
	if err != nil {
		saveErrToFile(err)
		log.Println(err)
	}
	lastData.LastSuccessfulFullScanStart = int64(lastSuccessfulFullScanStart)

	lastSuccessfulUpdate, _, err := key.GetIntegerValue(values.LAST_SUCCESSFUL_UPDATE_SUBKEY_NAME)
	if err != nil {
		saveErrToFile(err)
		log.Println(err)
	}
	lastData.LastSuccessfulUpdate = int64(lastSuccessfulUpdate)

	lastUpdate, _, err := key.GetIntegerValue(values.LAST_UPDATE_SUBKEY_NAME)
	if err != nil {
		saveErrToFile(err)
		log.Println(err)
	}
	lastData.LastUpdate = int64(lastUpdate)
	return
}

func stringToInt64(str string) (int64, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		saveErrToFile(err)
		return 0, err
	}
	return i, nil
}

func getTicket() (ticket Ticket, err error) {

	startPos := "{\"ActivationCenterId"
	endPos := "\"}"

	bs, err := getBinaryKeyInfoByKeyItem(values.REG_TICKET, values.TICKET_SUB_KEY)
	if err != nil {
		saveErrToFile(err)
		return
	}
	str := string(bs)
	is := strings.Index(str, startPos)
	ie := strings.Index(str, endPos)
	str = str[is : ie+len(endPos)]
	b := []byte(str)
	err = json.Unmarshal(b, &ticket)
	if err != nil {
		saveErrToFile(err)
		return
	}
	return
}

func getBinaryKeyInfoByKeyItem(keyPath string, keyItem string) (bytes []byte, err error) {
	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
	if err != nil {
		saveErrToFile(err)
		return nil, fmt.Errorf("registry.OpenKey uninstall err: %w", err)
	}
	defer regKey.Close()
	//keyInfo,err:= regKey.Stat()
	bytes, _, err = regKey.GetBinaryValue(keyItem)
	return
}

func getTimeFromDate(date string) (t time.Time, err error) {
	t, err = time.Parse(time.RFC3339, date)
	if err != nil {
		saveErrToFile(err)
		return t, err
	}
	return t, err
}

func getDateFromTimeStampSec(stamp int64) (date string) {
	date = time.Unix(stamp, 0).String()
	return
}

func getKesInfoValues() (values KesInfoValues, err error) {

	kesRoot, err := getKesRoot()

	if err != nil {
		saveErrToFile(err)
		return
	}
	values.REG_INSTALL_DATE_PATH = `SOFTWARE\WOW6432Node\KasperskyLab\protected\` + kesRoot + `\environment`
	values.REG_DATA = `SOFTWARE\WOW6432Node\KasperskyLab\protected\` + kesRoot + `\Data`
	values.REG_TICKET = `SOFTWARE\WOW6432Node\KasperskyLab\protected\` + kesRoot + `\watchdog\Ticket`
	values.TICKET_SUB_KEY = ""
	values.LAST_START_SUBKEY_NAME = `LastStart`
	values.LAST_SUCCESSFUL_FULL_SCAN_SUBKEY_NAME = `LastSuccessfulFullScan`
	values.LAST_SUCCESSFUL_FULL_SCAN_START_SUBKEY_NAME = `LastSuccessfulFullScanStart`
	values.LAST_SUCCESSFUL_UPDATE_SUBKEY_NAME = `LastSuccessfulUpdate`
	values.LAST_UPDATE_SUBKEY_NAME = `LastUpdate`
	return
}

func getKesRoot() (r string, err error) {
	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE, REG_CONNECTOR_CACHE, registry.READ)
	if err != nil {
		saveErrToFile(err)
		return "", fmt.Errorf("registry.OpenKey uninstall err: %w", err)
	}
	defer regKey.Close()
	rs, err := regKey.ReadSubKeyNames(1)
	if err != nil {
		saveErrToFile(err)
		return
	}
	if len(rs) > 0 {
		r = rs[0]
	}
	if len(r) == 0 {
		return "", fmt.Errorf("noRootVal: %w", err)
	}

	return
}

func saveErrToFile(saveErr error) {
	fileName := "./errors.txt"

	isHas, err := checkFileToHas(fileName)
	if !isHas {
		emptyFile, errCreate := os.Create(fileName)
		if errCreate != nil {
			//errCollector(err, "func (f *FileWorker)SaveProductPositionsInCategory(positions marketStructs.PositionsInCategory)(err error)")
			return
		}
		defer emptyFile.Close()
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {

		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(saveErr.Error()); err != nil {
		panic(err)
	}
}

func saveStructToFile(someStruct interface{}, fileName string) (err error) {
	arrayBytes, errJsonMarshal := json.Marshal(someStruct)

	isHas, err := checkFileToHas(fileName)

	if !isHas {
		emptyFile, errCreate := os.Create(fileName)
		if errCreate != nil {
			//errCollector(err, "func (f *FileWorker)SaveProductPositionsInCategory(positions marketStructs.PositionsInCategory)(err error)")
			return
		}
		defer emptyFile.Close()
	}

	if errJsonMarshal != nil {
		//errCollector(errJsonMarshal, "func (f *FileWorker)SaveProductPositionsInCategory(positions marketStructs.PositionsInCategory)(err error)")
		return
	}
	err = ioutil.WriteFile(fileName, arrayBytes, 0644)
	return

}

func checkFileToHas(fileURL string) (isHas bool, err error) {
	_, err = os.Stat(fileURL)
	if err != nil {

		isHas = false
		return
	}
	isHas = true
	return
}
