package pkg

import (
	"encoding/json"
	"log"
	"os"
)

type OneDumpFileJsonV1 struct {
	CrashModuleName    string
	CrashModuleVersion string
	ExceptionAddress   string
	ExceptionCode      string
	ExeModuleName      string
	ExeModuleVersion   string
	FilePath           string
	ModuleList         []string
	ModuleNum          int
	Version            string
}

type ModuleList struct {
	MoudleAndOffset string
	Percent         float64
	Number          int
}

type MoudleAndOffsetList struct {
	List []ModuleList

	Version string

	AllNumber    int
	ChromeNumber int
	Divide       float64
}

type OneVersionResultJson struct {
	MoudleAndOffset string
	Number          int
	Percent         float64
	SrcFile         []string
}

type OneVersionResultJsonList struct {
	List    []OneVersionResultJson
	Version string
}

func ReadOneDumpFileJson(filename string) (re OneDumpFileJsonV1) {

	file, e := os.OpenFile(filename, os.O_RDONLY, 0)
	if e != nil {
		log.Println(e.Error())
	}
	defer file.Close()
	enc := json.NewDecoder(file)

	err := enc.Decode(&re)
	if err != nil {
		log.Println("Error in Decode OneDumpFileJsonV1.json")
		log.Panicln(err.Error())
	}

	return re
}

func ReadMoudleAndOffsetListJson(filename string) (re MoudleAndOffsetList) {

	file, e := os.OpenFile(filename, os.O_RDONLY, 0)
	if e != nil {
		log.Println(e.Error())
	}
	defer file.Close()
	enc := json.NewDecoder(file)

	err := enc.Decode(&re)
	if err != nil {
		log.Println("Error in Decode MoudleAndOffsetList.json")
		log.Panicln(err.Error())
	}

	return re
}

func ReadJsonFile(filename string, re interface{}) error {

	file, e := os.OpenFile(filename, os.O_RDONLY, 0)
	if e != nil {
		log.Println(e.Error())
	}
	defer file.Close()
	enc := json.NewDecoder(file)

	e = enc.Decode(&re)
	if e != nil {
		log.Println("Error in Decode json", filename)
		log.Panicln(e.Error())
	}

	return e
}
