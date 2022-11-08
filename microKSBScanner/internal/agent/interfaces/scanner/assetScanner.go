package scanner

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/jaypipes/ghw"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/klauspost/cpuid"
	"github.com/matishsiao/goInfo"
	"github.com/shirou/gopsutil/v3/mem"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs/asset"
)

const s = `a$#2k2@$%#`

type AssetScanner struct {
}

var _ usecases.Scanner = &AssetScanner{}

func New() *AssetScanner {
	return &AssetScanner{}
}

func (a *AssetScanner) ID() (usecases.ID, error) {

	hostName, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("os.Hostname err: %w", err)
	}

	id, err := machineid.ID()
	if err != nil {
		return "", fmt.Errorf("mahoneid.ID err: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(s))
	mac.Write([]byte(id))
	mac.Write([]byte(hostName))
	return usecases.ID(hex.EncodeToString(mac.Sum(nil))), nil

}

func (a *AssetScanner) Scan() (asset.Asset, error) {
	hostid, err := a.ID()
	if err != nil {
		fmt.Errorf("a.ID err: %w", err)
		hostid = "not installed"
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Errorf("os.Hostname err: %w", err)
		hostname = "not determined"
	}

	macs, err := getMacAddr()
	if err != nil {
		fmt.Errorf("getMacAddr err: %w", err)
		var emptyMasc []string
		macs = emptyMasc
	}

	ips, err := getIpAddr()
	if err != nil {
		fmt.Errorf("getIpAddr err: %w", err)
		var emptyIps []string
		ips = emptyIps
	}

	gi, _ := goInfo.GetInfo()

	memory, err := mem.VirtualMemory()
	if err != nil {
		fmt.Errorf("mem.VirtualMemory err: %w", err)
	}

	block, err := ghw.Block()
	if err != nil {

		fmt.Errorf("ghw.Block err: %w", err)
	}

	hdds := make([]asset.HDDType, 0, len(block.Disks))

	for _, disk := range block.Disks {
		hdds = append(hdds, asset.HDDType{
			Name:         disk.Model,
			Manufacturer: disk.Vendor,
			Size:         bToGbString(disk.SizeBytes),
		})
	}

	gpu, err := ghw.GPU()
	var videos []string
	if err != nil {
		fmt.Errorf("ghw.GPU err: %w", err)
	} else {

		for _, card := range gpu.GraphicsCards {
			videos = append(videos, card.DeviceInfo.Vendor.Name+" "+card.DeviceInfo.Product.Name)
		}
	}
	var motherboard string
	mb, err := ghw.Baseboard()
	if err != nil {
		motherboard = "not determined"
	} else {
		motherboard = mb.String()
	}

	softwares, err := getSoftwareList()
	if err != nil {
		var emptySoftwares []asset.SoftwareType
		softwares = emptySoftwares
		fmt.Errorf("getSoftwareList err: %w", err)
	}

	return asset.Asset{
		HostID:     string(hostid),
		HostName:   hostname,
		Domain:     "",
		IPAddress:  strings.Join(ips, ","),
		MACAddress: strings.Join(macs, ","),
		Name:       hostname,
		OSName:     gi.OS,
		OSVersion:  gi.Core,
		CPU: []asset.CPUType{
			{
				Name: cpuid.CPU.BrandName, //+ " " + strconv.FormatInt(cpuid.CPU.Hz, 10),
			},
		},
		SystemMemory:  bToGbString(memory.Total),
		Video:         "", //strings.Join(videos, ","),
		HDD:           hdds,
		FactoryNumber: "",
		Motherboard:   motherboard,
		Monitor:       "",
		OD:            "",
		Software:      softwares,

		//Users: ,
	}, nil
}

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("net.Interfaces err: %w", err)
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func getIpAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("net.Interfaces err: %w", err)
	}
	var as []string
	for _, ifa := range ifas {
		addrs, err := ifa.Addrs()
		if err != nil {
			return nil, fmt.Errorf("ifa.Addrs err: %w", err)
		}
		for _, addr := range addrs {
			a := addr.String()
			if a != "" {
				as = append(as, a)
			}
		}
	}
	return as, nil
}

// здесь 1000 потому что единица измерения СИ GB, а есть ещё GiB которая с децтва называлась GB но значения были степени двойки 1024
// и похоже все производители измеряют в GB так как если измерять в GiB то не получаются цифры из прайс листа
func bToGbString(b uint64) string {
	n := b / 1000 / 1000 / 1000
	return strconv.FormatUint(n, 10) + " GB"
}
