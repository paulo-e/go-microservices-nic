package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-microservices-nic/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	p.l.Println("server: Returning list of products...")
}

func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	prod, _, err := data.GetProduct(id)
	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(rw).Encode(prod)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	p.l.Printf("server: Retrieving product...")
}

func (p *Products) PostProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	data.AddProduct(prod)
	p.l.Printf("server: Saving product...")
}

func (p *Products) PutProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	prod := &data.Product{}
	err = prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Product could not found", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
	p.l.Printf("server: Updating product...")
}

func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	err = data.DeleteProduct(id)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Product could not found", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
	p.l.Printf("server: Deleting product...")
}
