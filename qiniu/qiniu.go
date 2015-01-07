package qiniu

import (

	//"github.com/gogits/gogs/modules/mahonia"
	gio "io"

	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/rsf"
	//"golang.org/x/text/encoding"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Blog struct {
	ID      string
	Title   string
	Content string
}

func init() {
	ACCESS_KEY = ""
	SECRET_KEY = ""

	// 缓存 map 用以搜索
}

func Get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	s := string(body[:])
	return s
}

func Put(b Blog) {
	id := "t/" + b.ID + "/" + b.Title
	st := strings.NewReader(b.Content)
	uptoken := uptoken("proc")
	uploadBufDemo(st, id, uptoken)
}

func Query() {

}

//
func uptoken(bucketName string) string {
	putPolicy := rs.PutPolicy{
		Scope: bucketName,
		//CallbackUrl: callbackUrl,
		//CallbackBody:callbackBody,
		//ReturnUrl:   returnUrl,
		//ReturnBody:  returnBody,
		//AsyncOps:    asyncOps,
		//EndUser:     endUser,
		//Expires:     expires,
	}
	return putPolicy.Token(nil)
}

func uploadBufDemo(r gio.Reader, key, uptoken string) {
	// @gist uploadBuf
	var err error
	var ret io.PutRet
	var extra = &io.PutExtra{
	// Params:   params,
	// MimeType: mieType,
	// Crc32:    crc32,
	// CheckCrc: CheckCrc,
	}

	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器端生成的上传口令
	// key       为文件存储的标识
	// r         为io.Reader类型，用于从其读取数据
	// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	err = io.Put(nil, &ret, uptoken, key, r, extra)

	if err != nil {
		//上传产生错误
		log.Print("io.Put failed:", err)
		return
	}

	//上传成功，处理返回值
	log.Print(ret.Hash, ret.Key)
	// @endgist
}

func listAll(rs *rsf.Client, bucketName string, prefix string) {

	var entries []rsf.ListItem
	var marker = ""
	var err error
	var limit = 1000

	for err == nil {
		entries, marker, err = rs.ListPrefix(nil, bucketName, prefix, marker, limit)
		for _, item := range entries {
			//处理 item
			log.Print("item:", item)
		}
	}
	if err != gio.EOF {
		//非预期的错误
		log.Print("listAll failed:", err)
	}
}
