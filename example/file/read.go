// try.go
package main

import (
	//"bytes"
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		file, err := os.Open("Autumn.jpg") // For read access.
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Add("Content-Disposition", `filename=Aut柏木由纪.jpg`)
		w.Write(bytes)
	})
	fmt.Println(`http.ListenAndServe(":8080", nil)`)
	http.ListenAndServe(":8080", nil)
}
