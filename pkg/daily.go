package pkg

import (
	"fmt"
	"strings"
	"time"
)

func ConfigCycleAdjust(ini Config) (out Config) {
	out = ini
	if ini.CycleTimeType == "none" {
		out.DumpInputPath += out.CycleTimeDumpInputPathAppend
	} else {
		//计算今天，昨天，现在小时，过去一小时
		var internal_day int64 = 1e9 * 3600 * 24
		var internal_hour int64 = 1e9 * 3600
		now := time.Now()
		last_day := now.Add(-time.Duration(internal_day))
		last_hour := now.Add(-time.Duration(internal_hour))

		last_day_string := fmt.Sprintf("%02d", last_day.Day())
		now_day_string := fmt.Sprintf("%02d", now.Day())
		last_hour_string := fmt.Sprintf("%02d", last_hour.Hour())
		//now_hour_string := fmt.Sprintf("%02d", now.Hour())

		if ini.CycleTimeType == "day" {
			out.DumpInputPath = ini.DumpInputPath + "\\" + last_day_string
			out.DumpOutputPath = ini.DumpOutputPath + "\\" + last_day_string

			out.JsonInputPath = out.DumpOutputPath
			out.JsonOutputPath = ini.JsonOutputPath + "\\" + last_day_string

			out.TarInputPath = out.JsonOutputPath
			index := strings.LastIndex(ini.TarOutputPath, "\\")
			out.TarOutputPath = ini.TarOutputPath[:index+1] + last_day_string + ini.TarOutputPath[index:]

			out.SendTarPakagePath = out.TarOutputPath

			index = strings.LastIndex(ini.ReceiveTarPakagePath, "\\")
			out.ReceiveTarPakagePath = ini.ReceiveTarPakagePath[:index+1] + last_day_string + ini.ReceiveTarPakagePath[index:]

			out.UntarInputPath = out.ReceiveTarPakagePath
			out.UntarOutputPath = ini.UntarOutputPath + "\\" + last_day_string

			out.GenCallStackInputPath = out.UntarOutputPath

			//----------------------------
			out.DumpInputPath += out.CycleTimeDumpInputPathAppend

		} else if ini.CycleTimeType == "hour" {
			out.DumpInputPath = ini.DumpInputPath + "\\" + now_day_string + "\\" + last_hour_string
			out.DumpOutputPath = ini.DumpOutputPath + "\\" + now_day_string + "\\" + last_hour_string

			out.JsonInputPath = out.DumpOutputPath
			out.JsonOutputPath = ini.JsonOutputPath + "\\" + now_day_string + "\\" + last_hour_string

			out.TarInputPath = out.JsonOutputPath
			index := strings.LastIndex(ini.TarOutputPath, "\\")
			out.TarOutputPath = ini.TarOutputPath[:index+1] + now_day_string + "\\" + last_hour_string + ini.TarOutputPath[index:]

			out.SendTarPakagePath = out.TarOutputPath

			index = strings.LastIndex(ini.ReceiveTarPakagePath, "\\")
			out.ReceiveTarPakagePath = ini.ReceiveTarPakagePath[:index+1] +
				now_day_string + "\\" + last_hour_string + ini.ReceiveTarPakagePath[index:]

			out.UntarInputPath = out.ReceiveTarPakagePath
			out.UntarOutputPath = ini.UntarOutputPath + "\\" + now_day_string + "\\" + last_hour_string

			out.GenCallStackInputPath = out.UntarOutputPath

			//----------------------------
			out.DumpInputPath += out.CycleTimeDumpInputPathAppend
		}
	}

	return
}
