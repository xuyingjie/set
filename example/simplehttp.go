package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    args := os.Args
    path := "."
    if len(args) > 1 {
      path = args[1]    // arg[0] 是命令路径
    }
    fmt.Println(path + " on port 3000")
    http.ListenAndServe(":3000", http.FileServer(http.Dir(path)))
    //    http.Handle("/", http.FileServer(http.Dir(".")))
    //    http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
    //       fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    //  })
    // http.ListenAndServe(":3000", nil)
}