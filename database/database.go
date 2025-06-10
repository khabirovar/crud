package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetBookByID(id int) (*Book, error) {
	var book Book
	err := d.db.QueryRow("SELECT * FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (d *Database) GetBooks() ([]Book, error) {
	rows, err := d.db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]Book, 0)

	for rows.Next() {
		book := Book{}
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (d *Database) AddBook(book Book) error {
	_, err := d.db.Exec("INSERT INTO books(title, author, rating) values($1, $2, $3)", 
		book.Title, book.Author, book.Rating)
	return err
}

func (d *Database) UpdateBookByID(book Book) error {
	_, err := d.db.Exec("UPDATE books SET title=$1, author=$2, rating=$3 WHERE id = $4",
		book.Title, book.Author, book.Rating, book.ID)
	return err
}