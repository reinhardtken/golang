package pkg

import (
	"fmt"
	"strings"
	"util"
)

func WorkOnce(ini Config) {

	//-dump
	if strings.Contains(ini.WorkParams, "AnalyzeDumpFile") {
		fmt.Println("begin AnalyzeDumpFile")
		util.PrepareDir(ini.DumpOutputPath)
		re := Analyze(
			ini.DumpInputPath, ini.DumpOutputPath, ini.MoveDumpFile, ini.AnalyzeLevel,
			ini.DumpNum)
		fmt.Println("the num of dumps: ", re)
		fmt.Println("end AnalyzeDumpFile")
	}

	//-json
	if strings.Contains(ini.WorkParams, "GenJsonFile") {
		fmt.Println("begin GenJsonFile")
		util.PrepareDir(ini.JsonOutputPath)
		//pkg.GenDumpJsonFileByDump(ini.JsonInputPath, ini.JsonOutputPath, ini.JsonOutputMaxNum)
		GenDumpJsonFileByJson(ini.JsonInputPath, ini.JsonOutputPath, ini.JsonOutputMaxNum,
			ini.JsonModuleFilter)
		fmt.Println("end GenJsonFile")
	}

	//-tar
	if strings.Contains(ini.WorkParams, "DoTarPakage") {
		fmt.Println("begin DoTarPakage")
		//pkg.util.PrepareDir(ini.TarOutputPath)
		TarJsonAndDmp(ini)
		fmt.Println("end DoTarPakage")
	}

	//-send tar
	if strings.Contains(ini.WorkParams, "SendTarPakage") {
		fmt.Println("begin SendTarPakage")
		TcpSendFile(ini)
		fmt.Println("end SendTarPakage")
	}

	//-recv tar
	if strings.Contains(ini.WorkParams, "ReceiveTarPakage") {
		fmt.Println("begin ReceiveTarPakage")
		ReceiveTarPakage(ini)
		fmt.Println("end ReceiveTarPakage")
	}

	//-untar
	if strings.Contains(ini.WorkParams, "UntarPakage") {
		fmt.Println("begin UntarPakage")
		UntarFile(ini.UntarInputPath, ini.UntarOutputPath)
		fmt.Println("end UntarPakage")
	}

	//-callstack
	if strings.Contains(ini.WorkParams, "GenCallStackFile") {
		fmt.Println("begin GenCallStackFile")
		GenAllDirCallStack(ini.GenCallStackInputPath, ini.GenCallStackCmd, ini.GenCallStackParams,
			ini.GenCallStackVersionFilter)
		fmt.Println("end GenCallStackFile")
	}

	//-ondeir2json
	if strings.Contains(ini.WorkParams, "OneDirDump2Json") {
		fmt.Println("begin OneDirDump2Json")
		OneDirDump2Json(ini.OneDirDump2JsonInPath, ini.OneDirDump2JsonOutPath)
		fmt.Println("end OneDirDump2Json")
	}

	//-CallStackDistributed
	if strings.Contains(ini.WorkParams, "CallStackDistributed") {
		fmt.Println("begin CallStackDistributed")
		CallStackDistributed(ini.CallStackDistributedInPath, ini.CallStackDistributedOutPath,
			ini.CallStackDistributedMaxLevel, ini.CallStackDistributedMinNumber)
		fmt.Println("end CallStackDistributed")
	}
}
