package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, rr *http.Request) {
		log.Println("Hello world !!")
		out, err := ioutil.ReadAll(rr.Body)
		if err != nil {
			http.Error(rw, "Ooops", http.StatusBadGateway)
			return
		}
		log.Printf("Data is %s", out)
		fmt.Fprintf(rw, "Hello %s\n", out)
	})
	http.ListenAndServe(":9090", nil)
}
