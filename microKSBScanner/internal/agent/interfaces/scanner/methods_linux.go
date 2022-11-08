package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

func getSoftwareList() ([]asset.SoftwareType, error) {
	softwares, err := collect()
	if err != nil {
		return nil, fmt.Errorf("collect err: %w", err)
	}

	var sorted = make([]string, 0, len(softwares))

	for k := range softwares {
		sorted = append(sorted, k)
	}

	sort.Strings(sorted)

	var res = make([]asset.SoftwareType, 0, len(softwares))
	for _, name := range sorted {
		p := softwares[name]
		res = append(res, asset.SoftwareType{
			Name:         p.name,
			Version:      p.version,
			Manufacturer: p.maintainer,
			Description:  p.desc,
		})
	}

	return res, nil
}

type soft struct {
	name       string
	version    string
	maintainer string
	desc       string
}

func collect() (map[string]soft, error) {
	var short = make(map[string]soft)

	cmd := exec.Command("dpkg-query", "-f", "${Package};${Version};${Maintainer};${binary:Summary}\n", "-W")
	b, err := cmd.Output()
	if err != nil {
		return short, fmt.Errorf("output err: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ";")
		name := parts[0]
		version := parts[1]
		maintainer := parts[2]
		description := parts[3]

		if strings.HasPrefix(name, "lib") {
			continue
		}

		check(name, version, maintainer, description, &short)
	}

	if err := scanner.Err(); err != nil {
		return short, fmt.Errorf("scanner err: %w", err)
	}

	return short, nil
}

func check(name, version, maintainer, description string, short *map[string]soft) {
	parts := strings.Split(name, "-")

	tmp := parts[0]
	if _, ok := (*short)[tmp]; ok {
		return
	}
	for i := 1; i < len(parts); i++ {
		tmp = tmp + "-" + parts[i]
		if _, ok := (*short)[tmp]; ok {
			return
		}
	}

	(*short)[tmp] = soft{
		name:       name,
		version:    version,
		maintainer: maintainer,
		desc:       description,
	}
}
