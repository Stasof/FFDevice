package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var url = "http://192.168.1.125:8898/"
var baseRequest = BaseRequest{
	SerialNumber: "SNMOMF9100407",
	CheckCode:    "b64d6bcd",
}

func setUrlSerialAndCheck(ip string, serial string, check string) {
	url = "http://" + ip + ":8898/"
	baseRequest.SerialNumber = serial
	baseRequest.CheckCode = check
}

func getDetail() DetailResponse {
	detailRequest := DetailRequest{
		BaseRequest: baseRequest,
	}
	var res DetailResponse
	SendPOST("detail", detailRequest, &res)
	return res
}

func getProduct() ProductResponse {
	productRequest := ProductRequest{
		BaseRequest: baseRequest,
	}
	var res ProductResponse
	SendPOST("product", productRequest, &res)
	return res
}

func setLightControlCmd(status bool) CodeMessageResponse {
	sts := "close"
	if status {
		sts = "open"
	}
	cr := ControlRequest{
		BaseRequest: baseRequest,
		Payload: LightControlCmd{
			Cmd:  "lightControl_cmd",
			Args: LightArgs{Status: sts},
		},
	}
	var res CodeMessageResponse
	SendPOST("control", cr, &res)
	return res
}

func setCirculateCtlCmd(status_in bool, status_ext bool) CodeMessageResponse {
	sts_in := "close"
	sts_ext := "close"
	if status_in {
		sts_in = "open"
	}
	if status_ext {
		sts_ext = "open"
	}
	cr := ControlRequest{
		BaseRequest: baseRequest,
		Payload: CirculateCtlCmd{
			Cmd:  "circulateCtl_cmd",
			Args: CirculateArgs{Internal: sts_in, External: sts_ext},
		},
	}
	var res CodeMessageResponse
	SendPOST("control", cr, &res)
	return res
}

func SendPOST(path string, data any, response any) bool {

	// Кодируем структуру в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка кодирования JSON:", err)
		return false
	}

	fmt.Println(url + path)
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
