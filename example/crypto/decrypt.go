package main

import (
	"bytes"
	// "encoding/base32"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	num := 0
	root := "ln"
	if args := os.Args; len(args) > 1 {
		root = args[1]
	}

	filepath.Walk(root, func(eachPath string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fileName := path.Base(eachPath)
		if strings.Index(fileName, ".gpg") != -1 {
			fileName = strings.TrimLeft(strings.TrimLeft(fileName, "_"), ".")
			fileName = strings.TrimRight(strings.TrimRight(fileName, "gpg"), ".")
			// sDec, err0 := base32.StdEncoding.DecodeString(fileName)
			// if err0 != nil {
			// fmt.Println(fileName+" Decode error:", err)
			// return err0
			// }

			decPath := path.Dir(eachPath) + "/" + fileName

			cmd := exec.Command("gpg2", "--yes", "-o", decPath, "-d", eachPath)
			var out bytes.Buffer
			cmd.Stdout = &out
			e := cmd.Run()
			if e != nil {
				log.Fatal(e)
			}
			fmt.Printf(out.String())

			num++
		}
		return nil
	})
	fmt.Printf("%d", num)
}
