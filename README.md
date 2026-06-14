# FFDevice

Сервер, который может общаться с принтером и выводить информацию на HTML страничку.
Может работать как самомстоятельно, так и внутри OrcaSlicer (вкладка Device, собственно для этого и создавался)
Проверялась работа на FlashForge 5m pro с прошивкой 5.1.4 (но должно работать начиная с 3.2.7)
Часть протокола обмена взята отсюда https://github.com/Parallel-7/flashforge-api-docs/blob/main/endpoints/endpoints_5m_3.2.7.yaml

![ffdevice screenshot](https://github.com/Stasof/FFDevice/blob/main/ffdevice_screenshot.png)

**Возможности**

- Управление печатью (пауза, продолжение, отмена)
- Управление светом, вентиляторами (внешний, внутренний)
- Просмотр видео с видеокамеры
- Просмотр списка файлов на принтере
- Печать файла из памяти принтера
- Просмотр системной информации

**Компиляция из исходников**
- Linux: go build -o ffdevice main.go ffprotocol.go ffclient.go server.go
- Windows: GOOS=windows GOARCH=amd64 go build -o ffdevice.exe main.go ffprotocol.go ffclient.go server.go

**Скачать скомпилированные файлы (Windows, Linux)**

https://github.com/Stasof/FFDevice/releases

**Запуск**
- В папке с программой запустите ffdevice (.exe для windows).
- Программа запустится в консоли, либо без консоли (в зависимости от ОС), но ее можно увидеть в диспетчере задач.
- Вся отсальная настройка производится в web-браузере или OrcaSlicer.

**Подключение**

В строке web-браузера или в OrcaSlicer в поле URL-адрес хоста введите:

http://localhost:8765?ip=192.168.1.111&serial=SNMOMF7777777&check=b77d7bcd

Замените следующие параметры на свои:
- IP - IP адрес принтера
- SERIAL - серийный номе (слева вверху экрана, либо в информации об устройстве)
- CHECK - Printer ID (Настройки->сеть->сетевой режим->только локальные сети. Подробнее: https://www.flashforge.com/a/docs/orca-flashforge/orca-flashforge-quick-start-guide#connect-via-lan-only-mode)
