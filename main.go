package main

import (
	"net/http"

	"github.com/andrew-hillier/literate-sniffle/database"
	"github.com/andrew-hillier/literate-sniffle/product"
	"github.com/andrew-hillier/literate-sniffle/receipt"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = ""

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
