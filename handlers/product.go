package handlers

import (
	"log"
	"net/http"

	"github.com/rajat965ng/microservices/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	if rr.Method == http.MethodGet {
		getProducts(rw, rr)
		return
	}

	rw.WriteHeader(http.StatusNotImplemented)
}

func getProducts(rw http.ResponseWriter, rr *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to parse product list", http.StatusInternalServerError)
	}
}
