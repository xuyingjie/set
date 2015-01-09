/*
？title顺序，在前端调整
？上传图片插入到textarea，文件名+=TimeDiff()
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
	//"bytes"
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
	Id      string
	Keyword string
}
type user struct {
	Name   string
	Passwd string
}
type token struct {
	Name string
	Uid  string
}

func main() {

	var set []blog
	var index []anchor
	var tokenStore []token
	key := ``

	c := oss.NewClient()
	cache(c, key, &set, &index)

	//
	http.HandleFunc("/index", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, JsonEncode(index))
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, req *http.Request) {
		var q query
		JsonDecode(req.Body, &q)
		s := querySet(&set, &q)

		io.WriteString(w, JsonEncode(s))
	})

	http.HandleFunc("/put", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" && auth(req, &tokenStore) {
			var b blog
			bytes := JsonDecode(req.Body, &b)
			enc, _ := AESEncrypt([]byte(key), bytes)
			if b.Title != "" {
				id := putObject(c, b, string(enc))
				a := updateCache(b, id, &set, &index)
				io.WriteString(w, JsonEncode(a))
			}
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			var u user
			JsonDecode(req.Body, &u)
			if u.Name != "" && u.Passwd != "" {
				if login(w, c, u, &tokenStore) {
					io.WriteString(w, "200")
				}
			}
		}
	})
	http.HandleFunc("/auth", func(w http.ResponseWriter, req *http.Request) {
		if auth(req, &tokenStore) {
			io.WriteString(w, "200")
		}
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
		if logout(req, &tokenStore) {
			io.WriteString(w, "200")
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
//					filename += TimeDiff()
//				} else {
//					filename = strings.Replace(filename, ".", TimeDiff()+".", 1) // ?.tar.xz
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

// 缓存全部oss数据
func cache(c *oss.Client, key string, set *[]blog, index *[]anchor) {
	objectList, err := c.GetBucket("dbmy", "t/", "", "", "")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range objectList.Contents {
		bytes, err := c.GetObject("/dbmy/"+v.Key, -1, -1)
		if err != nil {
			fmt.Println(err)
		}
		dec, _ := AESDecrypt([]byte(key), bytes)
		var b blog
		var a anchor
		json.Unmarshal(dec, &b)
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
		id = TimeDiff()
	}

	err := c.PutObjectFromString("/dbmy/t/"+id+".json", s)
	if err != nil {
		fmt.Println(err)
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

// 或者使用成熟的session库
func login(w http.ResponseWriter, c *oss.Client, u user, tokenStore *[]token) bool {

	bytes, err := c.GetObject("/dbmy/etc/"+u.Name, -1, -1)
	if err != nil {
		fmt.Println(err)
	}

	if Sha1sum(u.Passwd) == string(bytes) {
		expiration := time.Now()
		expiration = expiration.AddDate(1, 0, 0)

		name := Sha1sum(u.Name)
		uid := Sha1sum(TimeDiff() + "31&rsv_t=4e3ek")
		exist := false
		for k, v := range *tokenStore {
			if name == v.Name {
				(*tokenStore)[k].Uid = uid
				exist = true
			}
		}
		if !exist {
			*tokenStore = append(*tokenStore, token{name, uid})
		}

		cookie := http.Cookie{Name: name, Value: uid, Expires: expiration}
		http.SetCookie(w, &cookie)
		return true
	} else {
		return false
	}
}

func auth(req *http.Request, tokenStore *[]token) bool {
	for _, cookie := range req.Cookies() {
		for _, v := range *tokenStore {
			if cookie.Name == v.Name && cookie.Value == v.Uid {
				return true
			}
		}
	}
	return false
}

func logout(req *http.Request, tokenStore *[]token) bool {
	for _, cookie := range req.Cookies() {
		for k, v := range *tokenStore {
			if cookie.Name == v.Name {
				(*tokenStore)[k].Uid = Sha1sum(TimeDiff() + "31&rsv_t=4es3ek")
				return true
			}
		}
	}
	return false
}
