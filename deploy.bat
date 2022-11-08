set version=v1
echo %version%
echo Build armAgent
start.\microKSBScanner\cmd\KSBAgent\build.bat %version%

echo Build armScanner
start .\microKSBScanner\cmd\KSBScanner\build.bat %version%




