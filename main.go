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
		if data.Args == "false false" {
			setCirculateCtlCmd(false, false)
		}
	case "detail":
		res := getDetail()
		return res

	case "product":
		res := getProduct()
		return res
	case "files":
		res := getFiles()
		return res
	case "thumb":
		res := getFileThumb(data.Args)
		return res
	case "command":
		fmt.Println(data.Cmd, data.Args)
		res := setCommand(data.Args)
		return res
		//case "print":
		//fmt.Println(data.Cmd, data.Args)
		//res := setCommand(data.Args)
		//return res
	}

	fmt.Println(data.Cmd, data.Args)

	return nil
}
