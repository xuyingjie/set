/*
？title顺序，在前端调整
？上传图片插入到textarea，文件名+=timeDiff()
？用户验证	http.Cookie
*/

package main

import (
	"io"
	"strings"
	//. "./qiniu"
	"fmt"
	"net/http"
	//"os"
	"./oss"
	"encoding/json"
	"strconv"
	//"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

type B struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Id      string `json:"id"`
}
type T struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}
type S struct {
	Id      string `json:"id"`
	Keyword string `json:"keyword"`
}

func main() {

	c := oss.NewClient()

	//
	//
	//
	//
	listObject, err := c.GetBucket("dbmy", "t/", "", "", "")
	if err != nil {
		log.Fatalln(err)
	}
	var title []T
	var set []B
	//set := make(map[string]blog)
	for _, v := range listObject.Contents {
		bytes, err := c.GetObject("dbmy/"+v.Key, -1, -1)
		if err != nil {
			log.Fatalln(err)
		}
		var b B
		var t T
		json.Unmarshal(bytes, &b)
		//b.Id = "/dbmy/" + v.Key
		b.Id = strings.TrimRight(strings.TrimLeft(v.Key, "t/"), ".json") // other
		t.Title = b.Title
		t.Id = b.Id
		title = append(title, t)
		set = append(set, b)
		//set[v.Key] = b
	}
	//
	//
	//

	//
	//
	//
	http.HandleFunc("/list", func(w http.ResponseWriter, req *http.Request) {
		b, err := json.Marshal(title)
		if err != nil {
			fmt.Println("json err:", err)
		}
		fmt.Println(string(b))
		io.WriteString(w, string(b))
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, req *http.Request) {

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
		}
		var search S
		json.Unmarshal(body, &search)

		var s []B
		for _, v := range set {
			if search.Id == v.Id {
				s = append(s, v)
			} else if search.Keyword != "" {
				re := regexp.MustCompile(`(?i)` + search.Keyword)
				if re.MatchString(v.Title) || re.MatchString(v.Content) {
					s = append(s, v)
				}
			}
		}

		b, err := json.Marshal(s)
		if err != nil {
			fmt.Println("json err:", err)
		}
		io.WriteString(w, string(b))
	})

	http.HandleFunc("/put", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Println(err)
			}
			var b B
			json.Unmarshal(body, &b)

			if b.Title != "" {
				var id string
				oldId := b.Id
				if oldId != "" {
					id = b.Id
				} else {
					id = timeDiff()
				}
				err = c.PutObjectFromString("/dbmy/t/"+id+".json", string(body))
				if err != nil {
					log.Fatalln(err)
				}

				var t T
				b.Id = id
				t.Title = b.Title
				t.Id = b.Id
				if oldId != "" {
					for k, v := range title {
						if b.Id == v.Id {
							title[k] = t
							set[k] = b
						}
					}
				} else {
					title = append(title, t)
					set = append(set, b)
				}

				tjson, err := json.Marshal(t)
				if err != nil {
					fmt.Println("json err:", err)
				}
				io.WriteString(w, string(tjson))
			}
		}

	})

	//http.HandleFunc("/upload", upload)

	http.Handle("/", http.FileServer(http.Dir("pub")))

	fmt.Println(`http.ListenAndServe(":8080", nil)`)
	http.ListenAndServe(":8080", nil)
}

//func upload(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "POST" {
//		r.ParseMultipartForm(32 << 20)
//		//file, handler, err := r.FormFile("file")
//		files := r.MultipartForm.File["file"]
//		var s string

//		for _, handler := range files {
//			filename := handler.Filename
//			file, err := handler.Open()
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//			defer file.Close()

//			if _, err := os.Stat("./pub/upload/" + filename); err == nil {
//				if strings.Index(filename, ".") == -1 {
//					filename += timeDiff()
//				} else {
//					filename = strings.Replace(filename, ".", timeDiff()+".", 1) // ?.tar.xz
//				}
//			}

//			path := "./pub/upload/" + filename
//			f, err := os.Create(path)
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//			defer f.Close()
//			io.Copy(f, file)

//			fi := strings.ToLower(filename)
//			if strings.HasSuffix(fi, ".png") || strings.HasSuffix(fi, ".jpg") || strings.HasSuffix(fi, ".gif") {
//				s += `<img style="max-width:100%;" src="` + path + `">`
//			} else {
//				s += `<div style="background:#f5f6f7;border:1px dashed #C9C9C9;padding:5px 10px;"><a href="` + path + `">` + filename + `</a></div><br>`
//			}
//		}
//		fmt.Fprint(w, s)
//	}
//}

func timeDiff() string {
	now := time.Now()
	then := time.Date(2014, 05, 10, 0, 0, 0, 0, time.Local)
	diff := now.Sub(then).Nanoseconds()

	return strconv.FormatInt(diff, 10)
}
