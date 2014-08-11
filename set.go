/*
// alpha
-不用ajax.
-用户验证
-跳转
-html模板
-UI设计(用bootstrap. foundation更多小屏优化；这两个只是主题细节不一样.)
-ckeditor(tinymce自动格式过滤有Bug)
-数据存储, noindex string
-删除 修改, 通过time.string查询删除
-搜索(没有标签</tag/*tag*>.-->直接搜索.)
-搜索结果加亮, 忽略< >中内容.
-多关键词搜索, 在输入框输入 a|b 可以实现或搜索.
-favicon.ico
-图片. 缺点无法迁移.
-单个文章页面 title/(title)
-多页面(offset实现). 用首页显示20, 搜索*显示全部替代.
-备份数据, 备份代码.

// beta
-把与数据相关操作提取出来,方便移植.
-switch to boltdb
-remove martini  // 各种框架琳琅满目，所以什么都不用。
-不需要登录验证，如果需要可以用http.Cookie。
-上传文件，文件名timeDiff()。
-指定数据文件
?? /_del 
*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func main() {

	dbPath := "_dbmy"
	if args := os.Args; len(args) > 1 {
		dbPath = args[1]
	}
	Open(dbPath)

	http.HandleFunc("/", index)
	http.HandleFunc("/put", put)
	http.HandleFunc("/modify", modify)
	http.Handle("/pub/", http.StripPrefix("/pub/", http.FileServer(http.Dir("pub"))))
	http.HandleFunc("/upload", upload)

	fmt.Println(`http.ListenAndServe(":8080", nil)`)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")

	if r.Method == "POST" {

		switch search {
		case "n":
			g := Blog{
				Title:   "",
				Content: "",
				ID:      "",
			}
			t, _ := template.New("editor").Parse(editor)
			t.Execute(w, g)
		default:
			if len([]byte(search)) <= 2 && search != "*" {
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				http.Redirect(w, r, "/?search="+search, http.StatusFound)
			}
		}

	} else {

		var blog, results []Blog

		if len([]byte(search)) <= 2 && search != "*" {
			search = ""
			blog = Query(10)
		} else {
			blog = Query(10000000000)
		}

		if len([]byte(search)) > 2 {

			re := regexp.MustCompile(`(?i)` + search) // `(?i)` +, 区分大小写

			for _, p := range blog {
				title := p.Title
				content := p.Content

				// 包括<>里的内容, 如果不包含需要多出很多计算	// strings.Index strings.Replace 无法忽略大小写
				if re.MatchString(title) || re.MatchString(content) {
					replace := `<span style="background-color: #ffff00;">${0}</span>`

					newTitle := re.ReplaceAllString(title, replace)
					findstr := regexp.MustCompile(`<[^<>]+>|[^<>]+`).FindAllString(content, -1) // 分割html 成 []string
					for k, v := range findstr {
						// 					if !regexp.MustCompile(`<[^<>]+>`).MatchString(v) {
						if strings.Index(v, "<") == -1 {
							findstr[k] = re.ReplaceAllString(v, replace)
						}
					}
					newContent := strings.Join(findstr, "")

					b := Blog{
						Title:   newTitle,
						Content: newContent,
						ID:      p.ID,
					}
					results = append(results, b)
				}
			}
		} else {
			results = blog
		}

		tmpl := struct {
			Search string
			Blog   []Blog
		}{search, results}

		t, _ := template.New("root").Parse(root)
		t.Execute(w, tmpl)
	}
}

func put(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		id := r.FormValue("id")
		if id != "" {
			Delete([]byte(id))
		}

		b := Blog{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			ID:      timeDiff(),
		}
		Put(b)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func modify(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		id := r.FormValue("id")
		switch r.FormValue("button") {
		case "Trash":
			Delete([]byte(id))
			http.Redirect(w, r, "/", http.StatusFound)
		case "Modify":
			blog := Get([]byte(id))
			t, _ := template.New("editor").Parse(editor)
			t.Execute(w, blog)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		//file, handler, err := r.FormFile("file")
		files := r.MultipartForm.File["file"]
		var s string

		for _, handler := range files {
			filename := handler.Filename
			file, err := handler.Open()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			if _, err := os.Stat("./pub/upload/" + filename); err == nil {
				if strings.Index(filename, ".") == -1 {
					filename += timeDiff()
				} else {
					filename = strings.Replace(filename, ".", timeDiff()+".", 1) // ?.tar.xz
				}
			}

			path := "./pub/upload/" + filename
			f, err := os.Create(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)

			fi := strings.ToLower(filename)
			if strings.HasSuffix(fi, ".png") || strings.HasSuffix(fi, ".jpg") || strings.HasSuffix(fi, ".gif") {
				s += `<img style="max-width:100%;" src="` + path + `">`
			} else {
				s += `<div style="background:#f5f6f7;border:1px dashed #C9C9C9;padding:5px 10px;"><a href="` + path + `">` + filename + `</a></div><br>`
			}
		}
		fmt.Fprint(w, s)
	}
}

func timeDiff() string {
	now := time.Now()
	then := time.Date(2014, 05, 10, 0, 0, 0, 0, time.Local)
	diff := now.Sub(then).Nanoseconds()

	return strconv.FormatInt(diff, 10)
}
