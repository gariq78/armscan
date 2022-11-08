 ECHO ON
 cd /d "%~dp0
 set armPath=..\..\..\arm\armScanner
 set bName=armScanner.exe
 set testPath=..\test_exe\
 echo %armPath%
 echo %bName%
 echo %testPath%
 go build -o %armPath%\%bName% -ldflags "-s -w -X main.version=0.0.0.0 -X main.name=ARMScanner"
 copy data.bin %armPath%\data.bin
 copy id.bin %armPath%\id.bin
 copy settings.bin %armPath%\settings.bin
 cd ..\..\..\arm
 xcopy "./armScanner" %testPath%\%1%\"armScanner\" /s /e /h
 xcopy "./armScanner" %testPath%\"httpFileServer\armScanner\" /s /e /h
 exit






