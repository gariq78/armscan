package infrastructure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/service/usecases"
)

func GetAgentZipsFromDir(dir string) (usecases.AgentZips, error) {
	var res = make(usecases.AgentZips)

	f, err := os.Open(dir)
	if err != nil {
		return res, fmt.Errorf("open %q err: %w", dir, err)
	}
	defer f.Close()

	files, err := f.Readdirnames(-1)
	if err != nil {
		return res, fmt.Errorf("Readdirnames err: %w", err)
	}

	for _, file := range files {
		bytes, err := ioutil.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return nil, fmt.Errorf("readfile %q err: %w", file, err)
		}

		seg := strings.Split(file, "_")
		os := strings.Split(seg[2], ".")

		res[usecases.OS(os[0])] = usecases.AgentZip{
			Version: seg[1],
			Zip:     bytes,
		}
	}

	return res, nil
}
