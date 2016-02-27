package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/rif/cache2go"
	"github.com/shunde/QRCode-Auth/qr"
	"github.com/shunde/QRCode-Auth/uuid"
	"image/jpeg"
	"log"
	"net/http"
	"time"
)

type qrenCode struct {
	name string
	data []byte
}

var (
	cache   *cache2go.CacheTable
	urlBase string = "http://localhost:8080/login/"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		var uid []byte
		for {
			uid = uuid.NewUuid()
			if !cache.Exists(string(uid)) {
				break
			}
		}
		m := qr.Encode(urlBase + string(uid))
		var buf bytes.Buffer
		jpeg.Encode(&buf, m, nil)

		var value qrenCode
		value.name = string(uid)
		value.data = buf.Bytes()

		cache.Cache(string(uid), 5*time.Minute, &value)
		w.Write(uid)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	}
}

func qrcode(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		qrcodeName := r.URL.Path[len("/qrcode/"):]
		res, err := cache.Value(qrcodeName)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 NotFound"))
			return
		}

		w.Header().Add("Content-Type", "image/jpeg")
		w.Write(res.Data().(*qrenCode).data)
	}
}

func main() {
	cache = cache2go.Cache("table")
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/qrcode/{name}", qrcode)
	r.HandleFunc("/login/{uid}", login)
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
