package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenAllDirCallStack(dir, cmd, params, versions string) {
	var now_output_dir string
	ignore := true

	walk := func(path string, info os.FileInfo, err error) error {

		if info.IsDir() == false && strings.HasSuffix(path, ".dmp") {
			if ignore == false {
				GenOneDmpCallStack(path, now_output_dir, cmd, params)
			}
		} else if info.IsDir() == true && info.Name() == "dmp" {
			index := strings.LastIndex(path, "\\")
			parent := path[:index+1]

			version := path[:index]
			index = strings.LastIndex(version, "\\")
			version = version[index+1:] 
			if strings.Contains(versions, version) {
				fmt.Println("callstack with version: ", version)
				ignore = false
			} else {
				ignore = true
			}

			callstack_path := parent + "callstack"
			fmt.Println("CreateDir:", callstack_path)
			os.MkdirAll(callstack_path, 0777)
			now_output_dir = callstack_path
		}
		return err
	}

	err := filepath.Walk(dir, walk)
	if err != nil {
		fmt.Println("no error expected, found:", err)
	}

}

func GenOneDmpCallStack(dmp_file, output_dir, cmd, params string) {
	/*
			logname = os.path.basename(file_name)
		  output_dir = os.path.dirname(file_name)
		  print(logname)
		  print(output_dir)

		  index = logname.find(".");
		  logname = logname[:index]
		  global WINDBG_PARAM
		  cmd = "windbg.exe -y O:\\360chrome_symbols;srv*E:\\dump\\symbols*http://msdl.microsoft.com/download/symbols  -c \""
		  #cmd = "cdb.exe -y C:\\Users\\liuqingping\\Desktop\\°æ±¾\\symbols;srv*E:\\dump\\symbols*http://msdl.microsoft.com/download/symbols  -c \""
		  cmd += WINDBG_PARAM
		  cmd += "\" -z "
		  cmd += (file_name)
		  cmd += ("     -logo ")
		  global OUTPUT_DIR
		  if not (len(OUTPUT_DIR) == 0) :
		    cmd += OUTPUT_DIR
		  else :
		    cmd += output_dir

		  cmd += "\\"

		  cmd += logname
		  cmd += (".log")
		  print(cmd)
		  return cmd
	*/
	index := strings.LastIndex(dmp_file, "\\")
	filename := dmp_file[index+1:]
	index = strings.LastIndex(filename, ".")
	filename = filename[:index]

	output_file := fmt.Sprintf("%s\\%s.log", output_dir, filename)
	cmdline := fmt.Sprintf("%s%s%s%s%s%s", cmd, params, "\" -z ", dmp_file, "     -logo ", output_file)

	fmt.Println("the windbg cmd line: ", cmdline)
	//用c调用好了。。。
	ExecCommand(cmdline)
	/*final_cmd := exec.Command("C:\\Program Files (x86)\\Debugging Tools for Windows (x86)\\windbg.exe", cmdline)
	err := final_cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	*/
	//output_cmd, _ := exec.Command(cmdline).Output()
	//fmt.Println(string(output))
}
