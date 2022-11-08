package nod32

import (
	"bufio"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const expirDateFileURL = "./"
const expirFile = "license.lf"

const startFindSTR = "<ESET schemaLocation"
const endSFindTR = "</PRODUCT_LICENSE_FILE></ESET>"
const findSTR = "0000-00-00T00:00:00Z"
const expirationDate = "FREE_PRODUCT EXPIRATION_DATE=\""

type Nod32 struct {
}

func (n *Nod32) GetExpiration() string {

	file, err := os.Open(expirDateFileURL + expirFile)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(transform.NewReader(file, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()))
	var nStr string
	for scanner.Scan() {
		nStr = nStr + scanner.Text()
	}
	sI := strings.Index(nStr, expirationDate) + len(expirationDate)

	return expirationDateFromString(nStr[sI : sI+len(findSTR)])
}
func expirationDateFromString(dateString string) string {
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {

		return "undefined"
	}

	time1stamp := time.Now().UTC()
	time2stamp := date.UTC()

	days := time2stamp.Sub(time1stamp).Hours() / 24
	sd := strconv.Itoa(int(days))
	rs := sd + " (" + dateString + ")"

	return rs
}
