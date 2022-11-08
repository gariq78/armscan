package drewb

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//!!!!  Спросить у Игоря корректно ли полчаю timestamps.
// Алгоритм такой.
// 1. считываем из реестра путь к исполняемому файлу
// 2. Выделяем из него путь до файла c key (в нем находиться интерсущяя нас дата о сроке окончания лицензии timestamps)
// 3. Находим в этой директории все файлы с раширением .key ( по уму он должен быть один. Хотя...... Надо мнение Игоря)
// 4. В первом файле полученого массива находим в файле начало timestamps и конец.
//    В const ex,se, вот тут есть вопрос будет ли этот поиск корректным ( в конце будут два не нужных символа \n\r как вариант прочекать до не числа)
// 5. Берем текущую дату и путем не сложных вычичлений находим количество оставшихся дней.
// 6. Возращаем назад сринг от инта . Так проще
const (
	// timestamps место окончания срока за этим стрингом в искомом файле
	ex = "EX="
	// конец
	se     = "SE="
	avName = "Dr.Web Security Space"
)

type DrWeb struct {
}

const uninstallRegistryPath = `Software\Microsoft\Windows\CurrentVersion\Uninstall`

// fileURL путь где лежит файл с ключем. в нем найдем файл с раширением .key
func (n *DrWeb) GetExpiration(fileUrl string) (strDays string) {

	files := glob(fileUrl, func(s string) bool {
		return filepath.Ext(s) == ".key"
	})

	if len(files) > 0 {
		if len(files) == 1 {
			strDays = getExpirationDays(files[0])
		} else {
			strDays = getExpirationDays(getLastModifyFile(files))
		}
	}
	return
}

func (n *DrWeb) GetInstallPath() (installUrl string, err error) {

	uninstallKey, err := registry.OpenKey(registry.LOCAL_MACHINE, uninstallRegistryPath, registry.READ)
	if err != nil {
		return "", fmt.Errorf("registry.OpenKey uninstall err: %w", err)
	}
	defer uninstallKey.Close()

	info, err := uninstallKey.Stat()
	if err != nil {
		return "", fmt.Errorf("uninstallKey.Stat err: %w", err)
	}

	subkeynames, err := uninstallKey.ReadSubKeyNames(int(info.SubKeyCount))
	if err != nil {
		return "", fmt.Errorf("uninstallKey.ReadSubKeyNames err: %w", err)
	}

	for _, name := range subkeynames {
		sub, err := registry.OpenKey(registry.LOCAL_MACHINE, uninstallRegistryPath+`\`+name, registry.READ)
		if err != nil {
			return "", fmt.Errorf("Open [%s] registry key err: %w", name, err)
		}

		dispName, _, _ := sub.GetStringValue("DisplayName")
		if dispName == avName {
			installUrl, _, err = sub.GetStringValue("DisplayIcon")
			if err != nil {

			}
			arrFromStr := strings.Split(installUrl, `\`)

			if len(arrFromStr) > 0 {
				arrFromStr = arrFromStr[:len(arrFromStr)-1]
			}

			installUrl = strings.Join(arrFromStr[:], `\`)
			break
		}

		sub.Close()
	}

	return
}

func getLastModifyFile(files []string) (file string) {
	var t int64 = 0
	for _, f := range files {
		tt, err := statTime(f)
		if err != nil {

		}
		if t < tt {
			t = tt
			file = f
		}
	}

	return
}

func getExpirationDays(fileURL string) string {

	dat, _ := os.ReadFile(fileURL)
	startPos := strings.Index(string(dat), ex)
	endPos := strings.Index(string(dat), se)

	nDat := dat[startPos+len(ex) : endPos]
	ss := string(nDat)
	var expTimestampsStr string
	// прочекааем полученное до не числа
	for _, char := range ss {
		s := string(char)
		_, err := strconv.Atoi(s)
		if err != nil {
			break
		}
		expTimestampsStr += s
	}

	// получаем timestamps из файла
	i, err := strconv.ParseInt(expTimestampsStr, 10, 64)
	if err != nil {
		log.Println(err)
	}

	expiresTime := time.Unix(i, 0)
	currentTimeStamp := time.Now().UTC()
	days := expiresTime.Sub(currentTimeStamp).Hours() / 24
	strDays := strconv.Itoa(int(days))

	return strDays
}

func glob(root string, fn func(string) bool) []string {
	var files []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if fn(s) {
			files = append(files, s)
		}
		return nil
	})
	return files
}

func statTime(name string) (timeUnix int64, err error) {
	fi, err := os.Stat(name)
	if err != nil {
		return
	}
	timeUnix = fi.ModTime().Unix()

	return
}
