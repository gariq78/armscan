 ECHO ON
 cd /d "%~dp0
 set armPath=..\..\..\arm\armAgent
 set bName=armAgent.exe
 echo %armPath%
 echo %bName%
 go build -o %armPath%\%bName% -ldflags "-s -w -X main.version=0.0.0.0 -X main.name=ARMAgent"
 copy settings.bin %armPath%\settings.bin
 cd %armPath%
 IF EXIST "assets" (
        echo ok
       ) ELSE (
          MkDir assets
       )
IF EXIST "agentZip" (
        echo ok
       ) ELSE (
          MkDir agentZip
       )


 cd ..\..\microKSBScanner\cmd\KSBAgent
