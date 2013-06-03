package main

import (
	"pkg"
	"util"
)

func main() {
	//pkg.CollectUnitTest()

	pkg.InitConfig()

	ini := pkg.GetConfig()
	ini2 := pkg.ConfigCycleAdjust(ini)
	//fmt.Println(ini2)
	util.ReflectPrintln(&ini2)

	//pkg.OneDirDump2Json(ini2.OneDirDump2JsonInPath, ini2.OneDirDump2JsonOutPath)
	pkg.CallStackDistributed(ini2.CallStackDistributedInPath, ini2.CallStackDistributedOutPath,
		ini2.CallStackDistributedMaxLevel, ini2.CallStackDistributedMinNumber)
}
