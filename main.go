package main

import "fmt"

func main() {
	run()
}

func Init(ip string, serial string, check string) {
	setUrlSerialAndCheck(ip, serial, check)
}

func API(data BaseRESTRequest) any {
	switch data.Cmd {
	case "light":
		if data.Args == "true" {
			setLightControlCmd(true)
		} else {
			setLightControlCmd(false)
		}
	case "fan":
		if data.Args == "false true" {
			setCirculateCtlCmd(false, true)
		}
		if data.Args == "true false" {
			setCirculateCtlCmd(true, false)
		}
	case "detail":
		res := getDetail()
		return res

	case "product":
		res := getProduct()
		return res
	}

	fmt.Println(data.Cmd, data.Args)

	return nil
}
