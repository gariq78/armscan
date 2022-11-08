ECHO ON
cd /d "%~dp0
set installPath=c:\testArmScanner
echo %installPath%
:m1
IF EXIST "%installPath%" (
      echo ok
      ) ELSE (
       mkdir %installPath%
        goto m1
      )


copy armScanner.exe %installPath%\armScanner.exe
copy data.bin %installPath%\data.bin
copy id.bin %installPath%\id.bin
copy settings.bin %installPath%\settings.bin
copy uninstall.bat %installPath%\uninstall.bat


:m3
IF EXIST "%installPath%\armScanner.exe" IF EXIST "%installPath%\id.bin" IF EXIST "%installPath%\data.bin" IF EXIST "%installPath%\settings.bin" IF EXIST "%installPath%\uninstall.bat " (
       echo ok
      ) ELSE (
        goto m3
      )
cd %installPath%
start armScanner.exe install
