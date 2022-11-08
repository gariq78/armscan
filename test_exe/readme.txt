/**********для тестирования win 10 (win7 пока без add_info)**************/

НАСТРОЙКА:

В папке armAgent правим файл settings.bin:
// тут меняем ip адрес и порт, на ip тачки на котором будет запущен агент 
"ServiceAddress":"http://xxx.xxx.xxx.xxx:xxxx/agent/api/v1" 

В папке httpFileServer правим файл settings.bin:
//тут меняем ip адрес и порт, на ip тачки и порт который будет слушать httpServ
"http_ip": "xxx.xxx.xxx.xxx", "http_port": "xxxx",


ЗАПУСК:

Агент:
В папке armAgent запускаем armAgent.exe. 
Он будет принмать данные от сканеров (для удобства папке assets будут еще лежать  сканы netNames.json)

HTTPServer (для скачки сканера):
из папкт httpFileServer запускаем httpFileServ.exe
После запуска, сформируеться armScanner.zip в папке public, доступный для скачивания в браузере по "ip":"port" httpFileServer

Сканер:
Скачиваем armScaner.zip
Сохраняем куда либо armScanner.zip или сразу распакованную папку armScanner

Если есть желание проинсталить, запускаем батник ОТ ИМЕНИ АДМИНИСТРАТОРА (пути установки там можно подправить в батнике, сейчас с:\testArmScanner)
Можно запустить без инсталяции , тогда он не впихнеться в службы винды.

После запуска сканера, смотрим все asset-ы в агенте





