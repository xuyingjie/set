// list.go
// 生成　src/resource.js　文件

package main

import (
	// "flag"
	// "fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getFilelist(root string, file string) {

	var pathArr []string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		path = strings.TrimLeft(path, "/")
		pathArr = append(pathArr, path)

		return nil
	})

	resources := "var RES = {\n"
	g_resources := "};\nvar g_resources = [\n"
	for i, v := range pathArr {

		a := strings.Split(path.Base(v), ".")

		// resName := strings.Replace(path.Base(v), ".", "_", -1)
		resName := ""
		if a[len(a)-1] == "png" {
			resName = "pic" + strings.Title(a[0])
		} else if a[len(a)-1] == "mp3" {
			resName = "music" + strings.Title(a[0])
		} else {
			resName = a[len(a)-1] + strings.Title(a[0])
		}

		if i == len(pathArr)-1 {
			g_resources += " RES." + resName + "\n"
			resources += " " + resName + ": \"" + v + "\"\n"
		} else {
			g_resources += " RES." + resName + ",\n"
			resources += " " + resName + ": \"" + v + "\",\n"
		}
	}
	resources += g_resources + "];\n"

	d := []byte(resources)
	ioutil.WriteFile(file, d, 0777)
}

func main() {
	// flag.Parse()
	// root := flag.Arg(0)
	// file := flag.Arg(1)
	getFilelist("res", "src/resource.js")
}
