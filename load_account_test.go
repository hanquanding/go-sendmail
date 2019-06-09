package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func loadAccountConf(path string) map[string]string {
	confMap := make(map[string]string)

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(file)
	// 循环读取
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// 过滤两端空格
		str := strings.TrimSpace(string(line))

		// 等号（=）的位置，没有找到跳过
		index := strings.Index(str, "=")
		if index < 0 {
			continue
		}

		// 等号（=）左边的值
		name := strings.TrimSpace(str[:index])
		if len(name) == 0 {
			continue
		}

		// 等号（=）右边的值
		address := strings.TrimSpace(str[index+1:])
		if len(address) == 0 {
			continue
		}

		confMap[name] = address
	}
	return confMap

}

// go test -v  --run="TestLoadAccountInfo"
func TestLoadAccountInfo(t *testing.T) {
	config := loadAccountConf("conf/account.txt")
	fmt.Println(config)

	for name, address := range config {
		fmt.Println("Name:", name, "Address:", address)
	}
}
