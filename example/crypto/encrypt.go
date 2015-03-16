/* encrypt.go
0.加密成 _.*.gpg 格式
3.可指定文件
4.文件名base32化. 只是防止搜索引擎抓取. 不方便选择解密哪个文件.(base64包含/字符, 所以不适合作为文件名)
5.不需要openpgp.js, 因为只是为了加密备份.
*/

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
	args := os.Args
	root := "set"
	if len(args) > 2 {
		root = args[2]
	}

	filepath.Walk(root, func(eachPath string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fileName := path.Base(eachPath)
		if strings.Index(fileName, ".gpg") == -1 {
			// base64Name := base32.StdEncoding.EncodeToString([]byte(fileName))
			encPath := path.Dir(eachPath) + "/" + "_." + fileName + ".gpg"

			// exec.Command("mkdir", "-p", fileDir).Run()
			cmd := exec.Command("gpg2", "--yes", "-o", encPath, "-r", args[1], "-e", eachPath)
			var out bytes.Buffer
			cmd.Stdout = &out
			e := cmd.Run()
			if e != nil {
				fmt.Println(fileName)
				log.Fatal(e)
			}
			fmt.Printf(out.String())

			num++
		}
		return nil
	})
	fmt.Printf("%d", num)
}
