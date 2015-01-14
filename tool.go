package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"
)

func TimeDiff() string {
	now := time.Now()
	then := time.Date(2014, 05, 10, 0, 0, 0, 0, time.Local)
	diff := now.Sub(then).Nanoseconds()

	return strconv.FormatInt(diff, 10)
}

func Sha1sum(s string) string {
	data := []byte(s)
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func JsonDecode(reader io.Reader, v interface{}) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		fmt.Println("json err:", err)
	}
}

func JsonEncode(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return bytes
}
