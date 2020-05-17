package handlers

import (
	"encoding/json"
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

	products := data.GetProducts()
	encoder := json.NewEncoder(rw)
	err := encoder.Encode(products)
	if err != nil {
		http.Error(rw, "Unable to parse product list", http.StatusInternalServerError)
	}
	// product, err := json.Marshal(products)
	// if err != nil {
	// 	http.Error(rw, "Unable to parse product list", http.StatusInternalServerError)
	// }
	// rw.Write(product)
}
