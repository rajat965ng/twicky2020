package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (hello *Hello) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	hello.l.Println("Hello world !!")
	out, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		http.Error(rw, "Ooops", http.StatusBadGateway)
		return
	}
	log.Printf("Data is %s", out)
	fmt.Fprintf(rw, "Hello %s\n", out)
}
