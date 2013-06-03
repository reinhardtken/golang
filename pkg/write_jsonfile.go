package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"util"
)

func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)

}

type filenode struct {
	dumppath string
	filename string
	number   int
}

type filenodes []filenode

func (p filenodes) Len() int           { return len(p) }
func (p filenodes) Less(i, j int) bool { return p[i].number >= p[j].number }
func (p filenodes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func GenDumpJsonFileByDump(dirin, dirout string, max int) {
	//根据目录树中dump文件的数目多少排序，对top max生成json，适用于目录树是单版本的情况
	var visited map[string]filenode
	visited = make(map[string]filenode)

	walk := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			index := strings.LastIndex(path, "\\")
			new_path := path[:index]
			_, ok := visited[new_path]
			if !ok {
				fmt.Println("the file will be analyze: %s", path)

				index := strings.LastIndex(new_path, "\\")
				filename := new_path[index+1:]
				filename = filename + ".json"
				visited[new_path] = filenode{path, filename, 1}
				//OneDump2Json(path, outfile)
			} else {
				visited[new_path] = filenode{visited[new_path].dumppath, visited[new_path].filename, visited[new_path].number + 1}
			}
		}
		return err
	}

	err := filepath.Walk(dirin, walk)
	if err != nil {
		fmt.Println("no error expected, found: %s", err)
	}
	//fmt.Println("the visited")
	//fmt.Println(visited)
	//fmt.Println("the visited after")

	//gen json file
	filenodelist := make(filenodes, len(visited))
	i := 0
	for _, v := range visited {
		filenodelist[i] = v
		i++
	}

	sort.Sort(filenodelist)
	//fmt.Println(filenodelist)
	//fmt.Println("before .....")

	for i, v := range filenodelist {
		s := fmt.Sprintf("%s\\%03d_%d_%s", dirout, i, v.number, v.filename)
		//fmt.Println(v.dumppath)
		//fmt.Println(s)

		OneDump2Json(v.dumppath, s)
		if i >= max {
			break
		}
	}

}

func GenDumpJsonFileByJson(dirin, dirout string, max int, JsonModuleFilter string) {

	fmt.Print("GenDumpJsonFileByJson begin \r\n", dirin, dirout)
	var in_jsonfile map[string]OneVersionResultJsonList
	in_jsonfile = make(map[string]OneVersionResultJsonList)
	json_file_slice := make([]string, 20)

	walk := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false && strings.HasSuffix(path, ".json") {
			json_file_slice = append(json_file_slice, path)

			index := strings.LastIndex(path, "\\")
			parent := path[:index]
			fmt.Print(parent, "\r\n", dirin)
			if parent == dirin {
				filename := path[index+1:]
				file_part := strings.Split(filename, "_")
				fmt.Print("\r\n", file_part[0], file_part[len(file_part)-1], "\r\n")
				if file_part[len(file_part)-1] == "4golang.json" {
					var content OneVersionResultJsonList
					if ReadJsonFile(path, &content) == nil {
						in_jsonfile[file_part[0]] = content
					} else {
						fmt.Println("read file failed: %s", path)
					}
				}
			}
		}
		return err
	}

	err := filepath.Walk(dirin, walk)
	if err != nil {
		fmt.Println("no error expected, found: %s", err)
	}
	//fmt.Println("the visited")
	//fmt.Println(visited)
	//fmt.Println("the visited after")

	//gen json file
	fmt.Println(in_jsonfile)

	//copy to json dir
	s0 := fmt.Sprintf("%s\\%s\\", dirout, "total")
	util.PrepareDir(s0)
	for _, v := range json_file_slice {
		index := strings.LastIndex(v, "\\")
		filename := v[index+1:]
		dst := s0 + filename
		CopyFile(v, dst)
	}

	for k, v := range in_jsonfile {

		s := fmt.Sprintf("%s\\%s\\%s", dirout, k, "json")
		s2 := fmt.Sprintf("%s\\%s\\%s", dirout, k, "dmp")
		//fmt.Println(v.dumppath)
		//fmt.Println(s)
		fmt.Println(s)
		util.PrepareDir(s)
		util.PrepareDir(s2)

		for index, v2 := range v.List {
			module := v2.MoudleAndOffset[:strings.LastIndex(v2.MoudleAndOffset, "\\")]
			if len(JsonModuleFilter) != 0 && strings.Contains(JsonModuleFilter, module) == false {
				continue
			}

			offset := strings.Replace(v2.MoudleAndOffset, "\\", "_", -1)
			file_head := fmt.Sprintf("%s\\%03d_%d_%s", s, index, v2.Number, offset)
			file_head2 := fmt.Sprintf("%s\\%03d_%d_%s", s2, index, v2.Number, offset)
			for i, v3 := range v2.SrcFile {
				file := fmt.Sprintf("%s_%02d.json", file_head, i)
				file2 := fmt.Sprintf("%s_%02d.dmp", file_head2, i)
				OneDump2Json(v3, file)
				//callstack
				//copy file
				CopyFile(v3, file2)

			}

			if index >= max {
				break
			}
		}

	}

}

//一个目录不支持嵌套
func OneDirDump2Json(in, out string) {
	util.PrepareDir(out)

	walk := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false && strings.HasSuffix(path, ".dmp") {
			index := strings.LastIndex(path, "\\")
			filename := path[index+1 : len(path)-4]
			filename += ".json"
			outname := out
			outname += "\\" + filename
			fmt.Println(outname)
			OneDumpStack2Json(path, outname)
		}
		return nil
	}

	err := filepath.Walk(in, walk)
	if err != nil {
		fmt.Println("no error expected, found: %s", err)
	}
}

/////////////////////////////////////////////////////////////
type CallStackJson struct {
	CallStackList      []string
	CrashModuleName    string
	CrashModuleVersion string
	ExceptionAddress   string
	ExceptionCode      string
	ExeModuleName      string
	ExeModuleVersion   string
	FilePath           string
	ModuleNum          int
	Version            string
}

func CallStackDistributed(in, out string, max, min int) {
	fmt.Println("CallStackDistributed", in, out, max)
	type OneLevel map[string]int
	levels := make([]OneLevel, max)
	//-这也太ugly了吧。。。
	for i, _ := range levels {
		levels[i] = make(map[string]int)
	}
	json_content := make(map[string]OneLevel)

	var files []string
	//files = make([]string, 10)
	walk := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false && strings.HasSuffix(path, ".json") {
			//fmt.Println(path)
			files = append(files, path)
			//fmt.Println("CallStackDistributed 4")
			//fmt.Println(len(files))
			//fmt.Println("CallStackDistributed 5")
		}
		return nil
	}

	err := filepath.Walk(in, walk)
	if err != nil {
		fmt.Println("no error expected, found: %s", err)
	}

	fmt.Println("CallStackDistributed 2")
	fmt.Println(files)
	fmt.Println(len(files))
	fmt.Println("CallStackDistributed 3")

	for _, v := range files {
		//fmt.Println(v)
		var one_stack CallStackJson
		ReadJsonFile(v, &one_stack)
		for i := 0; i < util.Min(max, len(one_stack.CallStackList)); i++ {
			value, ok := levels[i][one_stack.CallStackList[i]]
			if ok {
				levels[i][one_stack.CallStackList[i]] = value + 1
			} else {
				levels[i][one_stack.CallStackList[i]] = 1
			}
		}
	}

	for i, v := range levels {
		var new_one map[string]int
		new_one = make(map[string]int)
		for k, v2 := range v {
			if v2 > min {
				new_one[k] = v2
			}
		}
		levels[i] = new_one
	}

	for i := 0; i < max; i++ {
		key := fmt.Sprintf("%02d", i+1)
		json_content[key] = levels[i]
	}

	fmt.Println(json_content)
	WriteJsonFile(out, json_content)

}

func WriteJsonFile(filename string, re interface{}) error {
	fmt.Println("WriteJsonFile", filename)
	//util.ReflectPrintln(&re)

	file, e := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0)
	if e != nil {
		fmt.Println(e.Error())
	}
	defer file.Close()

	//encoder := json.NewEncoder(file)
	//e = encoder.Encode(&re)

	var b []byte
	b, e = json.MarshalIndent(&re, "\n", " ")
	if e != nil {
		fmt.Println("Error in Decode MoudleAndOffsetList.json")
		fmt.Println(e.Error())
	}

	_, e = file.Write(b)
	return e
}
