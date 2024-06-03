package api

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"

	"sf-news/pkg/storage"

	"github.com/gorilla/mux"
)

type API struct {
	storage       storage.StorageIface
	router        *mux.Router
	fileserverDir string
	indexFile     string
}

// Конструктор объекта программного интерфейса, где fileserverDir -
// корневая дирректория статического файл-сервера приложения,
// а storage - объект хранилища
func New(fileserverDir string, storage storage.StorageIface) *API {
	a := &API{
		storage:       storage,
		router:        mux.NewRouter(),
		fileserverDir: strings.TrimRight(fileserverDir, "/"),
		indexFile:     path.Join(fileserverDir, "index.html"),
	}
	a.endpoints()
	return a
}

// Возвращает внутренний роутер
func (a *API) Router() *mux.Router {
	return a.router
}

func (a *API) endpoints() {
	a.router.HandleFunc("/news/{count}", a.handlePosts).Methods(http.MethodGet)
	a.registerStatic()
}

func (a *API) handlePosts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(mux.Vars(r)["count"])
	posts, err := a.storage.PopPosts(count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (a *API) registerStatic() {
	// https://stackoverflow.com/questions/15834278/serving-static-content-with-a-root-url-with-the-gorilla-toolkit
	a.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(a.fileserverDir))))
	a.router.HandleFunc("/", a.handleIndex).Methods(http.MethodGet)
}

func (a *API) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, a.indexFile)
}
