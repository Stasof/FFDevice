package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// http://localhost:8765/?ip=192.168.1.111&serial=SNMOMF777777&check=b77d7bcd

func GetBR(printer PrinterStruct) BaseRequest {
	return BaseRequest{SerialNumber: printer.Serial, CheckCode: printer.Check}
}

func getDetail(printer PrinterStruct) DetailResponse {
	detailRequest := DetailRequest{
		BaseRequest: GetBR(printer),
	}
	var res DetailResponse
	SendPOST(printer, "detail", detailRequest, &res)
	return res
}

func getProduct(printer PrinterStruct) ProductResponse {
	productRequest := ProductRequest{
		BaseRequest: GetBR(printer),
	}
	var res ProductResponse
	SendPOST(printer, "product", productRequest, &res)
	return res
}

func setLightControlCmd(printer PrinterStruct, status bool) CodeMessageResponse {
	sts := "close"
	if status {
		sts = "open"
	}
	cr := ControlRequest{
		BaseRequest: GetBR(printer),
		Payload: LightControlCmd{
			Cmd:  "lightControl_cmd",
			Args: LightArgs{Status: sts},
		},
	}
	var res CodeMessageResponse
	SendPOST(printer, "control", cr, &res)
	return res
}

func setCommand(printer PrinterStruct, command string) CodeMessageResponse {
	cr := ControlRequest{
		BaseRequest: GetBR(printer),
		Payload: JobCtlCmd{
			Cmd:  "jobCtl_cmd",
			Args: JobArgs{JobID: "", Action: command},
		},
	}

	var res CodeMessageResponse
	SendPOST(printer, "control", cr, &res)
	return res
}

func setCirculateCtlCmd(printer PrinterStruct, status_in bool, status_ext bool) CodeMessageResponse {
	sts_in := "close"
	sts_ext := "close"
	if status_in {
		sts_in = "open"
	}
	if status_ext {
		sts_ext = "open"
	}
	cr := ControlRequest{
		BaseRequest: GetBR(printer),
		Payload: CirculateCtlCmd{
			Cmd:  "circulateCtl_cmd",
			Args: CirculateArgs{Internal: sts_in, External: sts_ext},
		},
	}
	var res CodeMessageResponse
	SendPOST(printer, "control", cr, &res)
	return res
}

func getFiles(printer PrinterStruct) GcodeListResponse {
	gcodeListRequest := GcodeListRequest{
		BaseRequest: GetBR(printer),
	}
	var res GcodeListResponse
	SendPOST(printer, "gcodeList", gcodeListRequest, &res)
	return res
}

func getFileThumb(printer PrinterStruct, filename string) GcodeThumbResponse {
	gcodeThumbRequest := GcodeThumbRequest{
		BaseRequest: GetBR(printer),
		FileName:    filename,
	}
	var res GcodeThumbResponse
	SendPOST(printer, "gcodeThumb", gcodeThumbRequest, &res)
	return res
}

func SendPOST(printer PrinterStruct, path string, data any, response any) bool {

	// Кодируем структуру в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка кодирования JSON:", err)
		return false
	}

	url := "http://" + printer.IP + ":8898/"
	// Отправляем POST-запрос
	resp, err := http.Post(
		url+path,                  // URL для тестирования
		"application/json",        // Content-Type
		bytes.NewBuffer(jsonData), // Тело запроса
	)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return false
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка сервера: %s\n", resp.Status)
		return false
	}

	// Декодируем JSON-ответ
	//var response any
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return false
	}

	return true
}
