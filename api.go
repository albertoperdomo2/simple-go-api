package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Image struct {
	FileName string  `json:"fileName"`
	Author   string  `json:"author"`
	ID       string  `json:"id"`
	Size     float32 `json:"size"`
}

type imageHandlers struct {
	sync.Mutex
	store map[string]Image
}

type adminPortal struct {
	password string
}

func (h *imageHandlers) images(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed\n"))
		return
	}
}

func (h *imageHandlers) get(w http.ResponseWriter, r *http.Request) {
	images := make([]Image, len(h.store))
	h.Lock()
	i := 0
	for _, image := range h.store {
		images[i] = image
		i++
	}
	h.Unlock()
	jsonBytes, err := json.Marshal(images)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *imageHandlers) getImage(w http.ResponseWriter, r *http.Request) {
	elements := strings.Split(r.URL.String(), "/")
	if len(elements) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if elements[2] == "random" {
		h.getRandomImage(w, r)
		return
	}
	h.Lock()
	image, ok := h.store[elements[2]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.Unlock()
	jsonBytes, err := json.Marshal(image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *imageHandlers) getRandomImage(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.store))
	h.Lock()
	i := 0
	for id := range h.store {
		ids[i] = id
		i++
	}
	defer h.Unlock()
	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}
	w.Header().Add("location", fmt.Sprintf("/images/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *imageHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}
	var image Image
	err = json.Unmarshal(bodyBytes, &image)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	image.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[image.ID] = image
	defer h.Unlock()
}

func newimageHandlers() *imageHandlers {
	return &imageHandlers{
		store: map[string]Image{},
	}
}

func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env variable ADMIN_PASSWORD not set")
	}
	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pswd, ok := r.BasicAuth()
	if !ok || pswd != a.password || user != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized\n"))
		return
	}
	w.Write([]byte("<html><h1>Super secret admin portal</h1></html>"))
}

func main() {
	admin := newAdminPortal()
	imageHandlers := newimageHandlers()
	http.HandleFunc("/images", imageHandlers.images)
	http.HandleFunc("/images/", imageHandlers.getImage)
	http.HandleFunc("/admin", admin.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
