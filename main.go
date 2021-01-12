package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonas-p/go-shp"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

//This program will extract all shp file table and create a csv copy of it.

func main() {
	defaultRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	root := flag.String("dir", defaultRoot, "root Directory")
	flag.Parse()
	fmt.Println(*root)
	// read directory
	shps := findFile(*root, ".shp")
	for _, f := range shps {
		shp2csv(f)
	}
}

func findFile(root string, match string) (file []string) {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if strings.HasSuffix(info.Name(), match) {
			file = append(file, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Total shp file : ", len(file))
	return file
}

func shp2csv(path string) {
	shape, err := shp.Open(path)
	if err != nil {
		fmt.Println("err")
		log.Fatal(err)
		return
	}
	defer shape.Close()
	fmt.Println(strings.Replace(path, ".shp", ".csv", -1))
	file, err := os.Create(strings.Replace(path, ".shp", ".csv", -1))
	if err != nil {
		fmt.Println("err")
		return
	}
	defer file.Close()

	fields := shape.Fields()
	f, _ := Decodebig5(fmt.Sprintf("%s\n", fields))
	f = strings.Replace(f, " ", ",", -1)
	file.WriteString(fmt.Sprintf("%s\n", f))
	for shape.Next() {
		n, _ := shape.Shape()

		for k := range fields {
			val := shape.ReadAttribute(n, k)
			val, _ = Decodebig5(val)
			if k != 0 {
				file.WriteString(",")
			}
			file.WriteString(fmt.Sprintf("%s", val))
		}
		file.WriteString(fmt.Sprintf("\n"))
	}
}

//Decodebig5 to utf8
func Decodebig5(s string) (string, error) {
	I := bytes.NewReader([]byte(s))
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return "", e
	}
	return string(d[:]), nil
}
