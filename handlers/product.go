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
		p.l.Println("Invoking GET call")
		getProducts(rw, rr)
		return
	}
	if rr.Method == http.MethodPost {
		p.l.Println("Invoking POST call")
		addProduct(rw, rr)
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

func addProduct(rw http.ResponseWriter, rr *http.Request) {
	p := &data.Product{}
	err := p.FromJSON(rr.Body)
	if err != nil {
		http.Error(rw, "Unable to parse JSON.", http.StatusBadRequest)
	}
	plist := data.AddProduct(p)
	plist.ToJSON(rw)
}
