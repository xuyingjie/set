package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

func main() {

	pub := `-----BEGIN PUBLIC KEY-----                                                                                                                  
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmuTuoRwtVRHfTYg/RR8o                                                                            
NI4AMFtBnwZoRPLa1eVXz2GaqwSYUvBzJ6TE7Uk394LXG1+sBGJ86R8G9bXRpmEI
F/toDdbVtZxF2kOBBUijUtaKUct3vZ9HiejInb4UTji6QI9IwNY2WfmkTZcB+nhW
vHmAOFc2D2CRBrkMzii+bzsEfs7eyrwaZkRr4atz/bzj9rbBcFAkdjyfhUfR7VA3
uhbADyDw2ZvXLZUklNzk49jVW6vy88Ma4b9K/Z2wREjR8YayKcqKy+ZjikxnFxqo
mn5qR0YCC8Kn0AS68aJgv9T1izN6ixOysDMH/GwSeAEywa6COs3LZ+cakJuEIlsF
6QIDAQAB
-----END PUBLIC KEY-----`
	pri := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAmuTuoRwtVRHfTYg/RR8oNI4AMFtBnwZoRPLa1eVXz2GaqwSY
UvBzJ6TE7Uk394LXG1+sBGJ86R8G9bXRpmEIF/toDdbVtZxF2kOBBUijUtaKUct3
vZ9HiejInb4UTji6QI9IwNY2WfmkTZcB+nhWvHmAOFc2D2CRBrkMzii+bzsEfs7e
yrwaZkRr4atz/bzj9rbBcFAkdjyfhUfR7VA3uhbADyDw2ZvXLZUklNzk49jVW6vy
88Ma4b9K/Z2wREjR8YayKcqKy+ZjikxnFxqomn5qR0YCC8Kn0AS68aJgv9T1izN6
ixOysDMH/GwSeAEywa6COs3LZ+cakJuEIlsF6QIDAQABAoIBACtgaxdBHczR+Xfr
7NP9mrPFBHaGR8Ux9SdB31jBJisUSU0VN1casRTsBp3wwDSXNoga7cA4mIRsRPWw
Wiq+x5Z3uDcP778xKTqrZslFAz/ihs/LLBg1L8KGJxBNt8AEENydg6FFG+lsiL+V
x18OLjQPzoF4otIBl/uK/PQK4vGg5OG+CZzbO1Q4Ahn15kqqOfNPKHKaDpOgZAgl
mhmpXmy25SlaNFJKZxb68BGtSJUBJlejatw6JDEO3F7DRdIDn2FyDApBW1OwuKtA
KRYCbkvfweEuNNaRmXj14pru5ZWEDXDXohSFB9UqNKTaIdQAASktGwRIOmMxgIt3
P7W75kUCgYEAyer2wIMT7QavhDvNyE9BraEyz0KMDDGEXcgBRmNeFIg4mQe8FfNu
t+dMbAvzNFi3005xDVcEo3smLmdPlvf/Y4f7eWXL6HVwGwQQIPSCrtI/vOFxGZAi
vg11GTfwTS3hjcX9HE4pR4G1uyp07USTiRhHoSN8WsqU4dmwYTADsgsCgYEAxGGs
BYXCM+/ZHAL8uxZEImFFy0sQq6ePfVWkXENb6LiGGKFh2MIRJldXyONWUXqv69kt
TdBXErBGNW9Z9gSOfFUE9YA5+kT/ZMx4KjUfzbphDOCXvjGUOTtS3QmkdJA9fMO7
l81Pkc2A+KOd+dUm0LAGFjPFHvBtzT2N0UNxtFsCgYEApfKMFbAk4jsKaV1VRPmO
ewru3VROEX9o0EKeeaEVIz7JdUvcExZcupxIPMyddzoq6mmflF0eHNYLjTuvN95e
cQjDbwRwz34lQq7WKp+J//AgHjYSY/YH97bLtIw63NOGeqRr36WFW5WJLGg6bP5d
WuEvjYnCnEO+lNf6lAWII/0CgYA9FKQMk63zuYYt0EALcMGAcADlWlO1EEjxEtIs
YEcV+066Gnf0k2gCJOiI8yzF6wMMuF/+8+4hQfKUbC3u9zvaMBd6xIdD8HH/SBmY
By39LxtAhhqsbX9MzcbYOUeNec+mHrsaXCGDmAelTj60ljeccSNzhGarWNzOGXci
v3d+QQKBgAorZh0pntXmOygJi50YHKj2zcXmQ728FixHwwnhUCK+98ZgQdmptKtq
ACstyN9TaB1p+a5J7UMPXAG6FLHJ6VSrp3UxRi8AV3I8IPeLmfRb6ugrhkak98sG
Y3Rc09Dhh2JK8EAiLVoSKEuiTVyVLOhH8+r0OzITU+WI4zh9hdF7
-----END RSA PRIVATE KEY-----`

	in := "Dbmy.xyz 柏木由纪"

	enStr := RSAEncrypt(in, pub)
	fmt.Println(enStr)

	s := RSADecrypt(enStr, pri)
	fmt.Println(s)

}

func RSAEncrypt(in, public string) string {
	block, _ := pem.Decode([]byte(public))
	if block == nil {
		fmt.Println("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}
	pub := pubInterface.(*rsa.PublicKey)

	sha1 := sha1.New()
	rand := rand.Reader

	encrypted, err := rsa.EncryptOAEP(sha1, rand, pub, []byte(in), nil)
	if err != nil {
		fmt.Println(err)
	}
	return string(encrypted)
}
func RSADecrypt(in, private string) string {
	block, _ := pem.Decode([]byte(private))
	if block == nil {
		fmt.Println("private key error")
	}
	priInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}
	pri := priInterface.(*rsa.PrivateKey)

	sha1 := sha1.New()

	plain, err := rsa.DecryptOAEP(sha1, nil, pri, []byte(in), nil)
	if err != nil {
		log.Println("error: %s", err)
	}

	return string(plain)
}
