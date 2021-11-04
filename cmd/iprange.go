package main

import (
	"fmt"
	"github.com/XinRoom/iprange"
	"io"
	"os"
	"strings"
)

func main() {
	// Stdin
	var isStdin bool
	if ss, err := os.Stdin.Stat(); err == nil && (ss.Mode()&os.ModeCharDevice) == 0 {
		isStdin = true
	}

	if len(os.Args) <= 1 && !isStdin {
		exeP := strings.Split(os.Args[0], "\\")
		fmt.Printf("Gen Ip Set.\n"+
			"Usage: %s ipStr [ipStr/file] ...\n"+
			"IP format can :\n"+
			"\t1.1.1.1\n"+
			"\t1.1.1.1-2\n"+
			"\t1.1.1-2.0-1\n"+
			"\t1.1.1.1/30\n"+
			"\t2001::59:63\n"+
			"\t2001::59:63-89\n"+
			"\t...\n"+
			"in addition: Support multiple parameters, file and commas\n", exeP[len(exeP)-1])
		return
	}
	var ipStrList []string

	// Stdin
	if isStdin {
		fcb, _ := io.ReadAll(os.Stdin)
		fc := string(fcb)
		ipStrList = append(ipStrList, strings.Split(fc, "\n")...)
	} else {
		// 多命令行参数
		for _, v := range os.Args[1:] {
			// , 号分割
			_v := strings.Split(v, ",")
			for _, v2 := range _v {
				// 文件解析
				if fileExists(v2) {
					f, _ := os.Open(v2)
					fcb, _ := io.ReadAll(f)
					fc := string(fcb)
					ipStrList = append(ipStrList, strings.Split(fc, "\n")...)
				} else {
					ipStrList = append(ipStrList, v2)
				}
			}
		}
	}

	for _, v := range ipStrList {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		it, startIp, err := iprange.NewIter(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] %s is not ip!\n", v)
			continue
		}
		for nit := startIp; it.HasNext(); nit = it.Next() {
			fmt.Println(nit)
		}
	}
}

// fileExists 判断文件存在
func fileExists(path string) bool {
	_, err := os.Lstat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
