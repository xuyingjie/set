package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
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

func JsonDecode(reader io.Reader, v interface{}) []byte {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return bytes
}

func JsonEncode(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return string(bytes)
}

// https://github.com/Unknwon/com/blob/master/string.go
// AESEncrypt encrypts text and given key with AES.
func AESEncrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// AESDecrypt decrypts text and given key with AES.
func AESDecrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
