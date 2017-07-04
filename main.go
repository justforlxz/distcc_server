package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	command := "/home/kazushin/distcc_scan/target/debug/distcc_finder"
	params := []string{"10.0.0.0/16"}
	execCommand(command, params)
}

func execCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false
	}

	fu := func() {
		os.Exit(1)
	}

	time.AfterFunc(20*time.Minute, fu)

	cmd.Start()

	os.Remove("/home/kazushin/.distcc/hosts")
	os.Create("/home/kazushin/.distcc/hosts")
	f, err := os.OpenFile("/home/kazushin/.distcc/hosts", os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("scaning...")
		reader := bufio.NewReader(stdout)

		//实时循环读取输出流中的一行内容
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			fmt.Println(strings.Replace(line, "\n", "", -1) + ",lzo,cpp")
			io.WriteString(f, strings.Replace(line, "\n", "", -1)+",lzo,cpp \n")
		}
	}

	cmd.Wait()
	return true
}
