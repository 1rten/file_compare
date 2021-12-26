package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"github.com/gookit/color"
)

var src string
var dst string
var f_type string

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func compareFile(file string) {
	if strings.HasSuffix(file, f_type) {
		src_path := fmt.Sprintf("%s/%s", src, file)
		dst_path := fmt.Sprintf("%s/%s", dst, file)
		src_f, err := os.Open(src_path)
		dst_f, _ := os.Open(dst_path)
		if err != nil {
			log.Fatal(err)
		}
		defer src_f.Close()
		defer dst_f.Close()
		src_r := bufio.NewReader(src_f)
		dst_r := bufio.NewReader(dst_f)
		color.Info.Println(fmt.Sprintf("[%s] COMPARE [%s]", src_path, dst_path))
		line := 0
		for {
			line++
			src_line, err := readLine(src_r)
			dst_line, _ := readLine(dst_r)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if src_line != dst_line && line != 1 {
				color.Warn.Println(fmt.Sprintf("LINE:%d {%s}", line, dst_line))
			}
		}
	}
}

func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = bs
	}
	return string(line), err
}

func init() {
	flag.StringVar(&src, "s", "", "输入源文件路径")
	flag.StringVar(&dst, "d", "", "输入目标文件路径")
	flag.StringVar(&f_type, "t", "", "输入文件类型")
}

func main() {
	flag.Parse()
	if !(isFlagPassed("s") && isFlagPassed("d") && isFlagPassed("t")) {
		flag.Usage()
		return
	}
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		compareFile(f.Name())
	}
}
