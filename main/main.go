package main

import (
	"fmt"
	"log"
	"pkg"
	"time"
	"util"
)

func main() {
	fmt.Println("crash analyze begin to run!")
	//init
	pkg.InitConfig()
	ini := pkg.GetConfig()
	ini = pkg.ConfigCycleAdjust(ini)
	//fmt.Println("WorkParams:", ini.WorkParams)
	util.ReflectPrintln(&ini)
	//prepare
	//Prepare(ini)
	//test
	//Test(ini)
	//work
	Work(ini)

}

//work//////////////////////////
func Work(ini pkg.Config) {
	if ini.WorkMode == "Once" {
		pkg.WorkOnce(ini)
	} else if ini.WorkMode == "Service" {
		Service(ini)
	} else {
		fmt.Println("wrong mode!!!")
	}
}

type Config struct {
	path pkg.Config
	now  time.Time
	last time.Time
}

type ZeroMonth int

func (z ZeroMonth) String() (s string) {
	if z > 9 {
		s = fmt.Sprintf("%d", z)
	} else {
		s = fmt.Sprintf("0%d", z)
	}

	return
}

func test_config() {
	ini := pkg.GetConfig()
	//fmt.Println(ini.DumpPath)
	//fmt.Println(ini.OutputPath)
	fmt.Println(ini.SymbolPath)
	fmt.Println(ini.TimeInternal)
}

func Service(ini pkg.Config) {
	var internal int64 = ini.TimeInternal
	internal *= 1e9
	last_time := time.Now()
	time.Sleep(time.Duration(internal))
	var config Config
	config.path = ini

	for i := 0; i < 10; i++ {
		config.now = time.Now()
		config.last = last_time
		/*config.y = now.Year()
		config.m = now.Month()
		config.d = now.Day()
		config.h = now.Hour()
		config.minute = now.Minute()
		config.s = now.Second()
		*/
		go work_once(config)

		last_time = config.now
		time.Sleep(time.Duration(internal))
	}

}

func work_once(c Config) {
	log.Println("work once ...")
	log.Println(c.now.Year(), ZeroMonth(int(c.now.Month())), c.now.Day(), c.now.Hour(), c.now.Minute(), c.now.Second())
	log.Println(c.last.Year(), int(c.last.Month()), c.last.Day(), c.last.Hour(), c.last.Minute(), c.last.Second())
}

//test////////////////////////////////////////////

func Test(ini pkg.Config) {
	//test 
	//json := pkg.ReadOneDumpFileJson("L:\\workspace\\dump_analyze\\Debug\\result\\dmp.json")
	//fmt.Println(json)

	//json2 := pkg.ReadMoudleAndOffsetListJson("L:\\workspace\\dump_analyze\\Debug\\result\\all_module_offset_report_4golang.json")
	//fmt.Println(json2)

	/*var json pkg.OneDumpFileJsonV1
	var json2 pkg.MoudleAndOffsetList
	pkg.ReadJsonFile("L:\\workspace\\dump_analyze\\Debug\\result\\dmp.json",
		&json)
	pkg.ReadJsonFile("L:\\workspace\\dump_analyze\\Debug\\result\\all_module_offset_report_4golang.json",
		&json2)
	fmt.Println(json)
	fmt.Println(json2)*/

	//test_config()
	//value := pkg.Hello()

	//fmt.Println(value)

	//pkg.TarJsonAndDmp(ini)
}
