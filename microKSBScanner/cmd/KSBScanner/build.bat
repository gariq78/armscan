 ECHO ON
 cd /d "%~dp0
 set armPath=..\..\..\arm\armScanner
 set bName=armScanner.exe
 echo %armPath%
 echo %bName%
 go build -o %armPath%\%bName% -ldflags "-s -w -X main.version=0.0.0.0 -X main.name=ARMScanner"
 copy data.bin %armPath%\data.bin
 copy id.bin %armPath%\id.bin
 copy settings.bin %armPath%\settings.bin
