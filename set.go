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

type blog struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Id      string `json:"id"`
}
type anchor struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}
type query struct {
	Id      string `json:"id"`
	Keyword string `json:"keyword"`
}

func main() {

	var set []blog
	var index []anchor

	c := oss.NewClient()
	cache(c, &set, &index)

	//
	http.HandleFunc("/index", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, jsonEncode(index))
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, req *http.Request) {
		var q query
		jsonDecode(req.Body, &q)
		s := querySet(&set, &q)

		io.WriteString(w, jsonEncode(s))
	})

	http.HandleFunc("/put", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			var b blog
			s := jsonDecode(req.Body, &b)
			if b.Title != "" {
				id := putObject(c, b, s)
				a := updateCache(b, id, &set, &index)
				io.WriteString(w, jsonEncode(a))
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

func jsonDecode(reader io.Reader, v interface{}) string {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return string(bytes)
}

func jsonEncode(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return string(bytes)
}

// 缓存全部oss数据
func cache(c *oss.Client, set *[]blog, index *[]anchor) {
	objectList, err := c.GetBucket("dbmy", "t/", "", "", "")
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range objectList.Contents {
		bytes, err := c.GetObject("/dbmy/"+v.Key, -1, -1)
		if err != nil {
			log.Fatalln(err)
		}
		var b blog
		var a anchor
		json.Unmarshal(bytes, &b)
		b.Id = strings.TrimRight(strings.TrimLeft(v.Key, "t/"), ".json") // other
		a.Title, a.Id = b.Title, b.Id
		*index = append(*index, a)
		*set = append(*set, b)
	}
}

func putObject(c *oss.Client, b blog, s string) (id string) {
	if b.Id != "" {
		id = b.Id
	} else {
		id = timeDiff()
	}

	err := c.PutObjectFromString("/dbmy/t/"+id+".json", s)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func updateCache(b blog, id string, set *[]blog, index *[]anchor) (a anchor) {
	if b.Id != "" {
		a.Title, a.Id = b.Title, b.Id
		for k, v := range *index {
			if b.Id == v.Id {
				(*index)[k] = a
				(*set)[k] = b
			}
		}
	} else {
		b.Id = id
		a.Title, a.Id = b.Title, b.Id
		*index = append(*index, a)
		*set = append(*set, b)
	}
	return
}

func querySet(set *[]blog, q *query) (s []blog) {
	for _, v := range *set {
		if q.Id == v.Id {
			s = append(s, v)
		} else if q.Keyword != "" {
			re := regexp.MustCompile(`(?i)` + q.Keyword)
			if re.MatchString(v.Title) || re.MatchString(v.Content) {
				s = append(s, v)
			}
		}
	}
	return
}
