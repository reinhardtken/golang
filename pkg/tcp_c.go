package pkg

import (
	"fmt"
	"net"
	"os"
)

const (
	READFILE_SUCC_GOON   = iota
	READFILE_SUCC_OVER   = iota
	READFILE_FAIL        = iota
	READFILE_CALLER_STOP = iota
)

func ReadFileFactory(file_name string) (func(buf []byte, stop bool) (int, int), error) {
	fmt.Println("ReadFile:", file_name)
	file, e := os.OpenFile(file_name, os.O_RDONLY, 0)
	if e != nil {
		fmt.Println(e.Error())
		return nil, e
	}

	var index int64 = 0
	return func(buf []byte, stop bool) (re, n int) {
		fmt.Println("ReadFileSon")
		if stop == false {
			var err error
			n, err = file.ReadAt(buf, index)
			if n != len(buf) {
				fmt.Println("ReadFile  READFILE_SUCC_OVER", err, n)
				re = READFILE_SUCC_OVER
			} else if n == len(buf) && err == nil {
				fmt.Println("ReadFile READFILE_SUCC_GOON", n)
				index += int64(n)
				re = READFILE_SUCC_GOON
			} else {
				fmt.Println("ReadFile READFILE_FAIL", err, n)
				re = READFILE_FAIL
			}
		} else {
			fmt.Println("ReadFile READ_CALLER_STOP")
			n = 0
			re = READFILE_CALLER_STOP
		}

		//clean if neccess
		if re != READFILE_SUCC_GOON {
			fmt.Println("ReadFile close")
			file.Close()
		}
		return
	}, nil

}

func ReadFileAndDoSomething(file_name string, f func([]byte) bool) {
	fmt.Println("OutputFile")
	ReadFile, err := ReadFileFactory(file_name)
	stop := false
	for {
		if err == nil {
			buf := make([]byte, 409800)
			e, n := ReadFile(buf, stop)
			if e == READFILE_SUCC_GOON || e == READFILE_SUCC_OVER {
				//这个地方也是有问题的，没有理由认为可以忽略send的返回值，如果send没发出去，返回值就不是期待的发送长度，这样
				//发送内容就丢了
				stop = f(buf[:n])
			}
			if e != READFILE_SUCC_GOON {
				break
			}
		}
	}
}

func Output(file_name string) {
	ReadFileAndDoSomething(file_name, func(buf []byte) bool {
		fmt.Println(string(buf))
		return false
	})
}

func TcpSendFile(ini Config) {
	// open connection:
	conn, err := net.Dial("tcp", ini.SendTarPakageServerInfo)
	if err != nil {
		// No connection could be made because the target machine actively refused it.
		fmt.Println("Error dialing", err.Error())
		return // terminate program
	}

	stop := false

	ReadFileAndDoSomething(ini.SendTarPakagePath, func(buf []byte) bool {
		fmt.Println("Send content len", len(buf))
		n, e := conn.Write(buf)
		if e != nil {
			fmt.Println("Send failed", e.Error())
			stop = true
		}
		if n != len(buf) {
			fmt.Println("Send error n != len(buf)")
		}
		return stop
	})
}

/*
func Example() {
	// open connection:
	conn, err := net.Dial("tcp", "localhost:50000")
	if err != nil {
		// No connection could be made because the target machine actively refused it.
		fmt.Println("Error dialing", err.Error())
		return // terminate program
	}
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("First, what is your name?")
	clientName, _ := inputReader.ReadString('\n')
	// fmt.Printf("CLIENTNAME %s",clientName)
	trimmedClient := strings.Trim(clientName, "\r\n") // "\r\n" on Windows, "\n" on Linux
	// send info to server until Quit:
	for {
		fmt.Println("What to send to the server? Type Q to quit.")
		input, _ := inputReader.ReadString('\n')
		trimmedInput := strings.Trim(input, "\r\n")
		// fmt.Printf("input:--%s--",input)
		// fmt.Printf("trimmedInput:--%s--",trimmedInput)
		if trimmedInput == "Q" {
			return
		}
		_, err = conn.Write([]byte(trimmedClient + " says: " +
			trimmedInput))
	}
}*/

////////////////////////////////////
