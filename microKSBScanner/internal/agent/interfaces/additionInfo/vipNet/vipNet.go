package vipNet

import (
	"bufio"

	"golang.org/x/text/encoding/charmap"

	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	SHARED_FILE = "c:\\ProgramData\\InfoTeCS\\ViPNet Client\\shared.cfg"
	CLIENT_FILE = "c:\\ProgramData\\InfoTeCS\\"
)

type Info struct {
	Net  string `json:"net"`
	Node string `json:"node"`
}

func GetInfo() (info []Info, err error) {
	data, err := readFromFile()
	info, err = parseDataToStruct(data)
	info = addNodeName(info)

	return
}

func addNodeName(info []Info) []Info {
	for i, _ := range info {
		clientPath := info[i].Node[i : len(info[i].Node)-1]
		fileName := clientPath[4:]
		fileUrl := CLIENT_FILE + "\\" + info[i].Node + "\\" + "APA" + fileName + ".TXT"
		info[i].Node = getNetNodeName(fileUrl)

	}
	return info
}

func parseDataToStruct(data []string) (info []Info, err error) {
	for _, str := range data {
		s := strings.Split(str, "\\")
		netInfo, _ := getNetNClientN(s[len(s)-1])
		info = append(info, netInfo)
	}
	return
}

func getNetNodeName(f string) (name string) {

	str, _ := decodeASCIIFile(f)

	ss := strings.Split(str, " ")
	var nbs []byte

	for _, b := range []byte(ss[len(ss)-1]) {
		if b == 13 {
			break
		}
		nbs = append(nbs, b)
	}
	name = " " + ss[len(ss)-2] + " " + string(nbs)
	return

}

func decodeASCIIFile(fileUrl string) (s string, err error) {

	f, err := os.Open(fileUrl)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	content := string(b)

	decoder := charmap.CodePage866.NewDecoder()
	reader := decoder.Reader(strings.NewReader(content))
	b, err = ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	s = string(b)
	return

}

func getNetNClientN(str string) (netInfo Info, err error) {
	snet := str[:4]
	netN, err := strconv.ParseUint(snet, 16, 64)

	netInfo.Net = "N " + strconv.FormatInt(int64(netN), 10)
	netInfo.Node = str

	return
}

func readFromFile() (s []string, err error) {
	f, err := os.Open(SHARED_FILE)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var b []byte
	scanner.Buffer(b, 100)
	i := 0
	for scanner.Scan() {
		if i != 0 {
			s = append(s, scanner.Text())
		}
		i++
	}
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
