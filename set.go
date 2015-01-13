/*
？title顺序，在前端调整
*/

package main

import (
	"./crypt"
	"./oss"
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

func main() {

	var set []blog
	var index []anchor
	var tokenStore []token
	key := ``

	c := oss.NewClient()
	cache(c, key, &set, &index)

	//
	m := macaron.Classic()

	m.Post("/index", func() string {
		return JsonEncode(index)
	})

	m.Post("/get", func(ctx *macaron.Context) string {
		var q query
		JsonDecode(ctx.Req.Body().ReadCloser(), &q)
		s := querySet(&set, &q)

		return JsonEncode(s)
	})

	m.Post("/put", func(w http.ResponseWriter, req *http.Request) {
		if auth(req, &tokenStore) {
			var b blog
			bytes := JsonDecode(req.Body, &b)
			enc, _ := crypt.Encrypt(key, bytes)
			if b.Title != "" {
				id := putObject(c, b, string(enc))
				a := updateCache(b, id, &set, &index)
				io.WriteString(w, JsonEncode(a))
			}
		}
	})

	m.Post("/login", func(w http.ResponseWriter, req *http.Request) {
		var u user
		JsonDecode(req.Body, &u)
		if u.Name != "" && u.Passwd != "" {
			if login(w, c, u, &tokenStore) {
				io.WriteString(w, "200")
			}
		}
	})
	m.Post("/auth", func(w http.ResponseWriter, req *http.Request) {
		if auth(req, &tokenStore) {
			io.WriteString(w, "200")
		}
	})
	m.Post("/logout", func(w http.ResponseWriter, req *http.Request) {
		if logout(req, &tokenStore) {
			io.WriteString(w, "200")
		}
	})

	m.Post("/upload", func(w http.ResponseWriter, req *http.Request) {
		upload(c, w, req)
	})
	m.Get("/p/:name", func(ctx *macaron.Context) {
		getPic(c, ctx)
	})

	m.Use(macaron.Static("pub"))
	m.Run()
}

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
		dec, _ := crypt.Decrypt(key, bytes)
		var b blog
		var a anchor
		json.Unmarshal(dec, &b)
		b.Id = strings.TrimLeft(v.Key, "t/")
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

	err := c.PutObjectFromString("/dbmy/t/"+id, s)
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

//
func upload(c *oss.Client, w http.ResponseWriter, r *http.Request) {
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

			//err = c.PutObjectFromReader("/dbmy/p/"+filename, file)
			//if err != nil {
			//	fmt.Println(err)
			//}

			path := "./cache/" + filename
			f, err := os.Create(path)
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()
			io.Copy(f, file)

			err = c.PutObject("/dbmy/p/"+filename, path)
			if err != nil {
				fmt.Println(err)
			}

			s += "\n![](/p/" + filename + ")\n"
		}
		fmt.Fprint(w, s)
	}
}

func getPic(c *oss.Client, ctx *macaron.Context) {

	filename := ctx.Params(":name")
	fmt.Println("Name: " + filename)

	if _, err := os.Stat("./cache/" + filename); err != nil {
		bytes, err := c.GetObject("/dbmy/p/"+filename, -1, -1)

		err = ioutil.WriteFile("./cache/"+filename, bytes, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	file, err := os.Open("./cache/" + filename) // For read access.
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	//.Header().Add("Content-Disposition", "filename="+filename)
	ctx.Resp.Write(bytes)
}
