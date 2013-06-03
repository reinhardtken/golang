package util

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

func PrepareDir(path string) {
	os.RemoveAll(path)
	time.Sleep(2 * time.Second)
	fmt.Println("PrepareDir create dir ", path)
	e := os.MkdirAll(path, 0777)
	if e != nil {
		fmt.Println("create dir failed :", e)
	}
}

func ReflectPrintln(unknown interface{}) {
	s := reflect.ValueOf(unknown).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\r\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

func Min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
	return a
}
