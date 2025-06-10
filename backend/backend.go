package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/khabirovar/crud/database"
)

type Backend struct {
	db   *database.Database
	port string
}

func NewBackend(dsn, port string) (*Backend, error) {
	db, err := database.NewDatabase(dsn)
	return &Backend{db: db, port: ":" + port}, err
}

func (b *Backend) Run() {
	http.HandleFunc("/books/", b.getBookByID)
	http.HandleFunc("/books", b.getBooks)

	if err := http.ListenAndServe(b.port, nil); err != nil {
		log.Fatal(err)
	}
}

func (b *Backend) getBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorMsg(err))
		return
	}

	book, err := b.db.GetBookByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(errorMsg(err))
		return
	}

	resp, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (b *Backend) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := b.db.GetBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}

	resp, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func errorMsg(err error) []byte {
	return []byte(fmt.Sprintf("{\"error\": \"%v\"", err))
}
