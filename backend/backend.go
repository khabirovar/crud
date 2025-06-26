package backend

import (
	"encoding/json"
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
	http.HandleFunc("/books", loggingMiddleware(b.handleBooks))
	http.HandleFunc("/books/", loggingMiddleware(b.handleBooks))

	if err := http.ListenAndServe(b.port, nil); err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s - %s\n", r.Method,r.RemoteAddr, r.URL)
		next(w, r)
	}
}

func (b *Backend) handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if len(strings.TrimPrefix(r.URL.Path, "/books")) > 1 {
			b.getBookByID(w, r)
		} else {
			b.getBooks(w, r)
		}
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
		errorMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := b.db.GetBookByID(id)
	if err != nil {
		errorMsg(w, http.StatusNotFound, err.Error())
		return
	}

	resp, err := json.Marshal(book)
	if err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMsg(w, http.StatusOK, resp)
}

func (b *Backend) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := b.db.GetBooks()
	if err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(books)
	if err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMsg(w, http.StatusOK, resp)
}

func (b *Backend) addBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	var book database.Book
	if err = json.Unmarshal(reqBytes, &book); err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = b.db.AddBook(book); err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := []byte(fmt.Sprintf("%v", map[string]string{"status": "created"}))
	jsonMsg(w, http.StatusCreated, resp)
}

func (b *Backend) updateBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	var book database.Book
	if err = json.Unmarshal(reqBytes, &book); err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}
	if book.ID == 0 {
		errorMsg(w, http.StatusBadRequest, "id must be specified")
		return
	}

	err = b.db.UpdateBookByID(book)
	if err != nil {
		errorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := []byte(fmt.Sprintf("%v", map[string]string{"status": "accepted"}))
	jsonMsg(w, http.StatusAccepted, resp)
}

func (b *Backend) deleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	err = b.db.DeleteBookByID(id)
	if err != nil {
		errorMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := []byte(fmt.Sprintf("%v", map[string]string{"status": "accepted"}))
	jsonMsg(w, http.StatusOK, resp)
}

// Helper functions
func errorMsg(w http.ResponseWriter, code int, message string) {
	payload := map[string]string{"error": message}
	response, _ := json.Marshal(payload)
	jsonMsg(w, code, response)
}

func jsonMsg(w http.ResponseWriter, code int, response []byte) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
