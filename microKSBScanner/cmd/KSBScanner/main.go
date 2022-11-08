package main

import (
	"os"

	"ksb-dev.keysystems.local/intgrsrv/microService/infrastructure/winservice"
)

var version = "0.0.0"
var name = "avrScanner"

func main() {

	app := &application{
		appversion: version,
		appname:    name,
	}

	serv := winservice.New(app)

	serv.Run(os.Args)
}
