# FFDevice

Сервер, который может общаться с принтером и выводить информацию на HTML страничку.
Может работать как самомстоятельно, так и внутри OrcaSlicer (собственно для этого и создавался)

Часть протокола обмена взята отсюда https://github.com/Parallel-7/flashforge-api-docs/blob/main/endpoints/endpoints_5m_3.2.7.yaml

Компиляция
Linux: go build -o ffdevice main.go ffprotocol.go ffclient.go server.go

Windows: GOOS=windows go build -o ffdevice.exe main.go ffprotocol.go ffclient.go server.go


