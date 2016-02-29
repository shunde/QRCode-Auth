package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rif/cache2go"
	"github.com/shunde/QRCode-Auth/uuid"
	"github.com/shunde/rsc/qr"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type qrenCodeInfo struct {
	name   string
	data   []byte
	scan   chan bool
	auth   chan bool
	isScan bool
	scanBy string
	token  string
}

type tokenInfo struct {
	ID         string // user ID
	deviceID   string
	timestamp  time.Time
	expireTime time.Time
}

type userInfo struct {
	name   string
	avatar []byte
	token  string
}

var (
	cache      *cache2go.CacheTable // uuid ==> qrenCodeInfo
	urlBase    string               = "http://localhost:8080/l/"
	tokenCache *cache2go.CacheTable // token ==> tokenInfo
	userCache  *cache2go.CacheTable // user ==> token
)

func JsLogin(w http.ResponseWriter, r *http.Request) {
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

		var value qrenCodeInfo
		value.name = string(uid) + ".png"
		value.data = c.PNG()
		value.scan = make(chan bool)

		cache.Cache(string(uid), 5*time.Minute, &value)

		respTpl := "window.code=%d; window.uuid='%s';"
		resp := fmt.Sprintf(respTpl, 200, string(uid))

		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(resp))
		w.WriteHeader(http.StatusOK)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("login.html")
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not Implemented."))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(dat)
	w.Header().Set("Content-Length", string(len(dat)))
}

func ScanQRCode(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		uuid := r.URL.Path[len("/l/"):]
		userID := r.FormValue("userID")
		fmt.Printf("user request from #%s#\n", userID)
		fmt.Printf("userCache: %d\n", userCache.Count())
		userinfo, err := userCache.Value(strings.TrimSpace(userID))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid User!"))
			return
		}
		fmt.Printf("scan request for uuid=%s from user %s\n", uuid, userinfo.Data().(*userInfo).name)
		res, err := cache.Value(uuid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 NotFound"))
			return
		}
		if res.Data().(*qrenCodeInfo).isScan {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		res.Data().(*qrenCodeInfo).scan <- true
		res.Data().(*qrenCodeInfo).isScan = true
		if res.Data().(*qrenCodeInfo).auth == nil {
			res.Data().(*qrenCodeInfo).auth = make(chan bool)
			res.Data().(*qrenCodeInfo).scanBy = userID
		}
		w.Write([]byte("scan success!"))
		fmt.Println("scan success!")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		uuid := r.URL.Query().Get("uuid")
		tip := r.URL.Query().Get("tip")

		//fmt.Printf("uuid=%s\ntip=%s\n", uuid, tip)
		if uuid == "" || tip == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := cache.Value(uuid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if tip == "0" { // query scan status
			if !res.Data().(*qrenCodeInfo).isScan {
				select {
				case <-res.Data().(*qrenCodeInfo).scan:
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/javascript")
					w.Write([]byte("window.code=201;"))
					// todo: add scanby info and user avatar
				case <-time.After(20 * time.Second):
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/javascript")
					w.Write([]byte("window.code=408;"))
				}
			}
			return
		} else if tip == "1" { // query authorization status
			select {
			case <-res.Data().(*qrenCodeInfo).auth:
				// get token
				// del uuid in cache
				token := res.Data().(*qrenCodeInfo).token
				cache.Delete(uuid)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/javascript")
				resp := "window.code=201; window.token='%s'"
				resp = fmt.Sprintf(resp, token)
				w.Write([]byte(resp))

			case <-time.After(20 * time.Second):
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/javascript")
				w.Write([]byte("window.code=408;"))

			}
			return
		}
	}

}

func Qrcode(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		uid := r.URL.Path[len("/qrcode/"):]
		res, err := cache.Value(uid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 NotFound"))
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(res.Data().(*qrenCodeInfo).data)

	}
}

func initServer() {
	cache = cache2go.Cache("table")
	cache.SetAboutToDeleteItemCallback(func(e *cache2go.CacheItem) {
		close(e.Data().(*qrenCodeInfo).scan)
		ch := e.Data().(*qrenCodeInfo).auth
		if ch != nil {
			close(ch)
		}
	})

	tokenCache = cache2go.Cache("tokenCache")
	userCache = cache2go.Cache("userCache")
}
func predefineUser() {
	users := []*userInfo{
		&userInfo{name: "张三"},
		&userInfo{name: "李四"},
		&userInfo{name: "隔壁老王"},
	}

	for i := 0; i < len(users); i++ {
		userCache.Cache(strconv.Itoa(i), 0, users[i])
	}
}

func main() {
	initServer()
	predefineUser()

	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/qrcode/{name}", Qrcode)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/l/{uuid}", ScanQRCode)
	r.HandleFunc("/jslogin", JsLogin)
	r.PathPrefix("/res/").Handler(http.FileServer(http.Dir(".")))
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
