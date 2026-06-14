package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PrintArgs struct {
	FileName            string `json:"fileName"`
	LevelingBeforePrint bool   `json:"levelingBeforePrint"`
}

func main() {
	run()
}

func API(data BaseRESTRequest) any {

	switch data.Cmd {
	case "light":
		if data.Args == "true" {
			setLightControlCmd(data.Printer, true)
		} else {
			setLightControlCmd(data.Printer, false)
		}
	case "fan":
		if data.Args == "false true" {
			setCirculateCtlCmd(data.Printer, false, true)
		}
		if data.Args == "true false" {
			setCirculateCtlCmd(data.Printer, true, false)
		}
		if data.Args == "false false" {
			setCirculateCtlCmd(data.Printer, false, false)
		}
	case "detail":
		res := getDetail(data.Printer)
		return res

	case "product":
		res := getProduct(data.Printer)
		return res
	case "files":
		res := getFiles(data.Printer)
		return res
	case "thumb":
		res := getFileThumb(data.Printer, data.Args)
		return res
	case "command":
		fmt.Println(data.Cmd, data.Args)
		res := setCommand(data.Printer, data.Args)
		return res
	case "print":
		var printArgs PrintArgs
		reader := strings.NewReader(data.Args)
		json.NewDecoder(reader).Decode(&printArgs)
		res := setPrintCmd(data.Printer, printArgs.FileName, printArgs.LevelingBeforePrint)
		return res
	}

	fmt.Println(data.Cmd, data.Args)

	return nil
}
