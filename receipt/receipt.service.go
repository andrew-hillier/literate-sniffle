package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/andrew-hillier/literate-sniffle/cors"
)

const receiptsBasePath = "receipts"

func SetupRoutes(apiBasePath string) {
	handleReceipts := http.HandlerFunc(recieptsHandler)
	handleReceipt := http.HandlerFunc(recieptHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptsBasePath), cors.Middleware(handleReceipts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptsBasePath), cors.Middleware(handleReceipt))
}

func recieptsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		receiptList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(receiptList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		r.ParseMultipartForm(5 << 20) // 5 Mb
		file, handler, err := r.FormFile("receipt")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		defer f.Close()
		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func recieptHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptsBasePath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	fContentType := http.DetectContentType(fHeader)
	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fSize := strconv.FormatInt(stat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", fSize)
	file.Seek(0, 0)
	io.Copy(w, file)
}
