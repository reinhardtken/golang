package pkg

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	WorkDescription string
	WorkMode        string
	WorkParams      string

	DumpDescription string
	DumpInputPath   string
	DumpOutputPath  string
	MoveDumpFile    bool
	AnalyzeLevel    int
	DumpNum         int

	JsonDescription  string
	JsonInputPath    string
	JsonOutputPath   string
	JsonOutputMaxNum int
	JsonModuleFilter string

	TarDescription string
	TarInputPath   string
	TarOutputPath  string

	UntarDescription string
	UntarInputPath   string
	UntarOutputPath  string

	GenCallStackDescription   string
	GenCallStackInputPath     string
	GenCallStackOutputPath    string
	GenCallStackCmd           string
	GenCallStackParams        string
	GenCallStackVersionFilter string

	SendTarPakagePath       string
	SendTarPakageServerInfo string

	ReceiveTarPakagePath            string
	ReceiveTarPakagePortInfo        string
	ReceiveTarPakageAfterWorkParams string

	CycleTimeType                string
	CycleTimeDumpInputPathAppend string

	OneDirDump2JsonInPath  string
	OneDirDump2JsonOutPath string

	CallStackDistributedInPath    string
	CallStackDistributedOutPath   string
	CallStackDistributedMaxLevel  int
	CallStackDistributedMinNumber int

	SymbolPath string

	TimeInternal int64
}

var ini Config

func init() {
	//write_test()
	//test()
	//file, e := os.OpenFile("/N/360云盘/workspace/golang/crash_analyze/src/src/cmd/config.json", os.O_RDONLY, 0)

}

func InitConfig() {
	file_path, _ := exec.LookPath(os.Args[0])
	index := strings.LastIndex(file_path, "\\")
	new_path := file_path[:index+1]
	final_path := new_path + "config.json"

	//file, e := os.OpenFile("N:\\360云盘\\workspace\\golang\\crash_analyze\\src\\src\\cmd\\config.json", os.O_RDONLY, 0)
	file, e := os.OpenFile(final_path, os.O_RDONLY, 0)
	if e != nil {
		log.Println(e.Error())
	}
	defer file.Close()
	enc := json.NewDecoder(file)

	err := enc.Decode(&ini)
	if err != nil {
		log.Println("Error in Decode config.json")
		log.Panicln(err.Error())
	}
}

func GetConfig() Config {
	return ini
}

func test() {
	const jsonStream = `
		{"Name": "Ed", "Text": "Knock knock."}
		{"Name": "Sam", "Text": "Who's there?"}
		{"Name": "Ed", "Text": "Go fmt."}
		{"Name": "Sam", "Text": "Go fmt who?"}
		{"Name": "Ed", "Text": "Go fmt yourself!"}
	`
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	for {
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s: %s\n", m.Name, m.Text)
	}
}

func write_test() {
	var m Config
	//m.DumpPath = "hello"
	//m.OutputPath = "world"
	m.SymbolPath = "!!!"
	file, e := os.OpenFile("N:\\360云盘\\workspace\\golang\\crash_analyze\\src\\src\\cmd\\config2.json", os.O_CREATE|os.O_WRONLY, 0)
	if e != nil {
		log.Println(e.Error())
	}
	defer file.Close()
	enc := json.NewEncoder(file)

	err := enc.Encode(m)
	if err != nil {
		log.Println("Error in Encode config.json")
	}
}
