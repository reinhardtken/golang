package pkg

import (
	"C"
	"syscall"
	"unsafe"
)

func Hello() int {
	//return int(C.hello())
	return 0
}

var (
	dll32 = syscall.NewLazyDLL("dump_analyze_dll.dll")

	analyze               = dll32.NewProc("analyze")
	one_dump_2_json       = dll32.NewProc("one_dump_2_json")
	execute_cmd           = dll32.NewProc("ExecCommand")
	one_dump_stack_2_json = dll32.NewProc("one_dump_stack_2_json")
)

func OneDump2Json(input_path, output_path string) {
	one_dump_2_json.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(input_path))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(output_path))))
}

func OneDumpStack2Json(input_path, output_path string) {
	one_dump_stack_2_json.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(input_path))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(output_path))))
}

func ExecCommand(path string) {
	execute_cmd.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))))
}

func Analyze(input_path, output_path string, movefile bool, analyze_level, dump_num int) int {
	var move C.int
	if movefile == true {
		move = 1
	} else {
		move = 0
	}

	re, _, _ := analyze.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(input_path))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(output_path))),
		uintptr(move), uintptr(analyze_level), uintptr(dump_num))
	return int(re)
}

func CallCFunc(input_path, output_path string) int {
	dll32 := syscall.NewLazyDLL("dump_analyze_dll.dll")
	println("call dll:", dll32.Name)
	g := dll32.NewProc("test1")
	func2 := dll32.NewProc("analyze")
	//cstr := C.CString("string from golang")
	//defer C.free(unsafe.Pointer(cstr))
	g.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("string from golang"))))
	//func2.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("L:\\workspace\\dump_analyze\\Debug\\BAK"))),
	//	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("L:\\workspace\\result"))))
	re, _, _ := func2.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(input_path))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(output_path))))
	return int(re)
}
