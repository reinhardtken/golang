package pkg

import (
	"fmt"
	"net"
	"os"
)

const (
	WRITEFILE_SUCC_GOON   = iota
	WRITEFILE_FAIL        = iota
	WRITEFILE_CALLER_STOP = iota
)

type WriteFileError struct {
	value int
}

func (*WriteFileError) Error() string {
	return "WriteFileError"
}

func WriteFileFactory(file_name string) (func(buf []byte, stop bool) WriteFileError, error) {
	fmt.Println("WriteFile:", file_name)
	file, e := os.OpenFile(file_name, os.O_WRONLY|os.O_CREATE, 0)
	if e != nil {
		fmt.Println(e.Error())
		return nil, e
	}

	var index int64 = 0
	return func(buf []byte, stop bool) (re WriteFileError) {
		fmt.Println("WriteFileSon")
		if stop == false {
			var n int
			var err error
			n, err = file.WriteAt(buf, index)
			if n == len(buf) && err == nil {
				fmt.Println("WriteFile WRITEFILE_SUCC_GOON", n)
				index += int64(n)
				re.value = WRITEFILE_SUCC_GOON
			} else {
				fmt.Println("WriteFile WRITEFILE_FAIL", err, n)
				re.value = WRITEFILE_FAIL
			}
		} else {
			fmt.Println("WriteFile WRITEFILE_CALLER_STOP")
			re.value = WRITEFILE_CALLER_STOP
		}

		//clean if neccess
		if re.value != WRITEFILE_SUCC_GOON {
			fmt.Println("WriteFile close")
			file.Close()
		}
		return
	}, nil

}

func ReceiveTarPakage(ini Config) {
	fmt.Println("Starting the server ...")
	// create listener:
	listener, err := net.Listen("tcp", ini.ReceiveTarPakagePortInfo)
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return // terminate program
	}
	// listen and accept connections from clients:
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // terminate program
		}
		fmt.Println("one peer connectted")
		go GoWork(ini, conn)
	}
}

func WriteFileFromSource(filename string, Source func(bool, []byte) (error, int)) (e error) {
	WriteFile, _ := WriteFileFactory(filename)
	if WriteFile != nil {
		buf := make([]byte, 409800)
		for {
			//这个地方是有问题的，不能在socket 读到数据后去等磁盘io，如果磁盘io太久
			//tcp的buf就爆了，虽然这个概率很低
			var n int
			e, n = Source(false, buf)
			if e != nil {
				fmt.Println("Error Source", e.Error())
				WriteFile(buf, true)
				//认为socket关闭就是传输完成。。。
				return nil // terminate program 
			} else {
				e2 := WriteFile(buf[:n], false)
				if e2.value != WRITEFILE_SUCC_GOON {
					return &e2
				}
			}
		}
	}
	return nil
}

func GoWork(ini Config, conn net.Conn) {

	WriteFileFromSource(ini.ReceiveTarPakagePath, func(stop bool, buf []byte) (e error, n int) {
		if stop == false {
			n, e = conn.Read(buf)
			if e != nil {
				fmt.Println("Read failed", e.Error())
				//fmt.Println("read content", string(buf))
				return e, 0
			}
			return nil, n
		}

		return nil, 0
	})

	//if e == nil {
	//改动作参数,执行后续untar，callstack的动作
	fmt.Println("after receive tar file")
	fmt.Println(ini.ReceiveTarPakageAfterWorkParams)
	ini.WorkParams = ini.ReceiveTarPakageAfterWorkParams
	WorkOnce(ini)
	//}
}

func Example() {
	fmt.Println("Starting the server ...")
	// create listener:
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return // terminate program
	}
	// listen and accept connections from clients:
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // terminate program
		}
		go Work(conn)
	}
}

func Work(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return // terminate program
		}
		fmt.Printf("Received data: %v", string(buf))
	}
}
