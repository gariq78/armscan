module ksb-dev.keysystems.local/intgrsrv/microKSBScanner

go 1.16

replace ksb-dev.keysystems.local/intgrsrv/microService => ../microService

require (
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/go-chi/chi v1.5.4
	github.com/jaypipes/ghw v0.9.0
	github.com/klauspost/cpuid v1.3.1
	github.com/matishsiao/goInfo v0.0.0-20210923090445-da2e3fa8d45f
	github.com/mattn/go-sqlite3 v1.14.14
	github.com/mitchellh/go-ps v1.0.0
	github.com/shirou/gopsutil/v3 v3.22.6
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8
	golang.org/x/text v0.3.7
	ksb-dev.keysystems.local/intgrsrv/microService v0.0.0-00010101000000-000000000000
)
