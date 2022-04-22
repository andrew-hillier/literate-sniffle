package main

import (
	"net/http"

	"github.com/andrew-hillier/literate-sniffle/product"
)

const apiBasePath = ""

func main() {
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
