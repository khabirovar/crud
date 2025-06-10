package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	http.HandleFunc("/books", b.handleBooks)
	http.HandleFunc("/books/", b.handleBooks)

	if err := http.ListenAndServe(b.port, nil); err != nil {
		log.Fatal(err)
	}
}

func (b *Backend) handleBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("METHOD: %v", r.Method)
	switch r.Method {
	case http.MethodGet:
		if len(strings.TrimPrefix(r.URL.Path, "/books")) > 1 {
			b.getBookByID(w, r)
		}
		b.getBooks(w, r)
	case http.MethodPost:
		b.addBook(w, r)
	case http.MethodPatch:
		b.updateBook(w, r)
	case http.MethodDelete:
		b.deleteBook(w, r)
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

func (b *Backend) addBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}

	var book database.Book
	if err = json.Unmarshal(reqBytes, &book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}

	if err = b.db.AddBook(book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (b *Backend) updateBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}

	var book database.Book
	if err = json.Unmarshal(reqBytes, &book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}
	if book.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorMsg(errors.New("id must be specified")))
		return
	}

	err = b.db.UpdateBookByID(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg(err))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (b *Backend) deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HIER")
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	fmt.Printf("id %s", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorMsg(err))
		return
	}
	fmt.Printf("id = %d\n", id)

	err = b.db.DeleteBookByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorMsg(err))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// TODO transform to msg func with application json header and code
func errorMsg(err error) []byte {
	return []byte(fmt.Sprintf("{\"error\": \"%v\"", err))
}
