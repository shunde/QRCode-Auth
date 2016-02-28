package main

import (
	//"bytes"
	"github.com/gorilla/mux"
	"github.com/rif/cache2go"
	//"github.com/shunde/QRCode-Auth/qrcode"
	"github.com/shunde/QRCode-Auth/uuid"
	"github.com/shunde/rsc/qr"
	//"image/png"
	"log"
	"net/http"
	"time"
)

type qrenCodeInfo struct {
	name   string
	data   []byte
	isScan bool
}

type token struct {
}

var (
	cache   *cache2go.CacheTable
	urlBase string = "http://localhost:8080/l/"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		var uid []byte
		for {
			uid = uuid.NewUuid()
			if !cache.Exists(string(uid)) {
				break
			}
		}

		c, _ := qr.Encode(urlBase+string(uid), qr.H)

		var value qrenCode
		value.name = string(uid) + ".png"
		value.data = c.PNG()

		cache.Cache(string(uid), 5*time.Minute, &value)

		w.Write(uid)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	}
}

func Qrcode(w http.ResponseWriter, r *http.Request) {
	qrcodeName := r.URL.Path[len("/qrcode/"):]
	res, err := cache.Value(qrcodeName)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 NotFound"))
		return
	}

	if r.Method == "GET" {
		w.Header().Add("Content-Type", "image/png")
		w.Write(res.Data().(*qrenCodeInfo).data)
	} else if r.Method == "POST" {
		// mark qrcode scaned
		res.Data().(*qrenCodeInfo).isScan = true
	}
}

func main() {
	cache = cache2go.Cache("table")
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/qrcode/{name}", Qrcode)
	r.HandleFunc("/l/{uid}", Login)
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
