package handlers

import (
	"log"
	"net/http"
)

type GoodBye struct {
	l *log.Logger
}

func NewGoodBye(l *log.Logger) *GoodBye {
	return &GoodBye{l}
}

func (g *GoodBye) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	g.l.Println("Inside Goodbye handler..")
	rw.Write([]byte("GoodBye..."))
}
