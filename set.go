package main

import (
	"code.google.com/p/rsc/crypt"
	"encoding/json"
	"fmt"
	"github.com/xuyingjie/set/oss"
	"net/http"
	"regexp"
	"strings"
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

var set []blog
var index []anchor
var tokenStore []token
var c *oss.Client

const key = ``

func main() {

	c = oss.NewClient()
	cache()

	//
	http.HandleFunc("/index", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, string(JsonEncode(index)))
	})
	http.HandleFunc("/get", querySet)
	http.HandleFunc("/put", putObject)

	//
	http.HandleFunc("/login", login)
	http.HandleFunc("/auth", func(w http.ResponseWriter, req *http.Request) {
		if auth(w, req) {
			fmt.Fprintf(w, "true")
		}
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
		if logout(w, req) {
			fmt.Fprintf(w, "true")
		}
	})

	//
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/p/", getPic)

	http.Handle("/", http.FileServer(http.Dir("pub")))

	fmt.Println(`http.ListenAndServe(":8080", nil)`)
	http.ListenAndServe(":8080", nil)
}

// 缓存全部oss数据
func cache() {
	objectList, err := c.GetBucket("dbmy", "t/", "", "", "")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range objectList.Contents {
		bytes, err := c.GetObject("/dbmy/"+v.Key, -1, -1)
		if err != nil {
			fmt.Println(err)
		}
		dec, _ := crypt.Decrypt(key, bytes)
		var b blog
		var a anchor
		json.Unmarshal(dec, &b)
		b.Id = strings.TrimLeft(v.Key, "t/")
		a.Title, a.Id = b.Title, b.Id
		index = append(index, a)
		set = append(set, b)
	}
}

func querySet(w http.ResponseWriter, req *http.Request) {

	var q query
	var s []blog

	JsonDecode(req.Body, &q)

	for _, v := range set {
		if q.Id == v.Id {
			s = append(s, v)
		} else if q.Keyword != "" {
			re := regexp.MustCompile(`(?i)` + q.Keyword)
			if re.MatchString(v.Title) || re.MatchString(v.Content) {
				s = append(s, v)
			}
		}
	}
	fmt.Fprintf(w, string(JsonEncode(s)))
}

func putObject(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" && auth(w, req) {

		var b blog
		JsonDecode(req.Body, &b)

		if b.Title != "" {

			var a anchor
			var id string

			if b.Id != "" {
				id = b.Id
				a.Title, a.Id = b.Title, b.Id
				for k, v := range index {
					if b.Id == v.Id {
						index[k] = a
						set[k] = b
					}
				}
			} else {
				b.Id = TimeDiff()
				id = b.Id
				a.Title, a.Id = b.Title, b.Id
				index = append(index, a)
				set = append(set, b)
			}

			bytes := JsonEncode(b)
			enc, _ := crypt.Encrypt(key, bytes)
			err := c.PutObjectFromString("/dbmy/t/"+id, string(enc))
			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprintf(w, string(JsonEncode(a)))
		}
	}

}

// 或者使用成熟的session库
func login(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {

		var u user
		JsonDecode(req.Body, &u)

		if u.Name != "" && u.Passwd != "" {

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
				for k, v := range tokenStore {
					if name == v.Name {
						(tokenStore)[k].Uid = uid
						exist = true
					}
				}
				if !exist {
					tokenStore = append(tokenStore, token{name, uid})
				}

				cookie := http.Cookie{Name: name, Value: uid, Expires: expiration}
				http.SetCookie(w, &cookie)
				fmt.Fprintf(w, "true")
			}
		}
	}
}

func auth(w http.ResponseWriter, req *http.Request) bool {
	for _, cookie := range req.Cookies() {
		for _, v := range tokenStore {
			if cookie.Name == v.Name && cookie.Value == v.Uid {
				return true
			}
		}
	}
	return false
}

func logout(w http.ResponseWriter, req *http.Request) bool {
	for _, cookie := range req.Cookies() {
		for k, v := range tokenStore {
			if cookie.Name == v.Name {
				(tokenStore)[k].Uid = Sha1sum(TimeDiff() + "31&rsv_t=4es3ek")
				return true
			}
		}
	}
	return false
}

//
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		//file, handler, err := r.FormFile("file")
		files := r.MultipartForm.File["file"]
		var s string

		for _, v := range files {
			filename := TimeDiff() + "_" + v.Filename
			file, err := v.Open()
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			err = c.PutObjectFromReader("/dbmy/p/"+filename, file)
			if err != nil {
				fmt.Println(err)
			}

			s += "\n![](/p/" + filename + ")\n"
		}
		fmt.Fprint(w, s)
	}
}

func getPic(w http.ResponseWriter, req *http.Request) {

	path := req.URL.Path
	filename := strings.TrimLeft(path, "/p/")

	bytes, err := c.GetObject("/dbmy/p/"+filename, -1, -1)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Add("Content-Disposition", "filename="+filename)
	w.Write(bytes)
}
