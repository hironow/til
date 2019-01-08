package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"google.golang.org/api/iterator"
)

// Entity

type BookID string

func NewBookID() (BookID, error) {
	uID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return BookID(uID.String()), nil
}

func GenerateBookKey(user *User, id BookID) *datastore.Key {
	return datastore.NameKey("Book", string(id), GenerateUserKey(user.ID))
}

type Book struct {
	ID     BookID
	UserID UserID
	Name   string
}

func NewBook(user *User, name string) (*Book, error) {
	id, err := NewBookID()
	if err != nil {
		return nil, err
	}
	return &Book{ID: id, UserID: user.ID, Name: name}, nil
}

// Store

type BookStore struct {
	client *datastore.Client
	debug  bool
}

func NewBookStore(cfg *Config, client *datastore.Client) *BookStore {
	return &BookStore{client: client, debug: cfg.Debug}
}

func (b *BookStore) Get(ctx context.Context, user *User, id BookID) (*Book, error) {
	var book Book
	err := b.client.Get(ctx, GenerateBookKey(user, id), &book)

	if b.debug {
		log.Printf("Book: %#v", book)
	}

	return &book, err
}

func (b *BookStore) Put(ctx context.Context, user *User, book *Book) error {
	_, err := b.client.Put(ctx, GenerateBookKey(user, book.ID), book)
	return err
}

func (b *BookStore) Delete(ctx context.Context, user *User, id BookID) error {
	return b.client.Delete(ctx, GenerateBookKey(user, id))
}

func (b *BookStore) List(ctx context.Context, user *User, cursorString string) ([]*Book, error) {
	q := datastore.NewQuery("Book").Ancestor(GenerateUserKey(user.ID))
	if cursorString != "" {
		cursor, err := datastore.DecodeCursor(cursorString)
		if err != nil {
			return nil, err
		}
		q.Start(cursor)
	}

	var books []*Book
	t := b.client.Run(ctx, q)
	for {
		var book Book
		_, err := t.Next(&book)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if b.debug {
		for i, book := range books {
			log.Printf("#%d Book: %#v", i, book)
		}
	}

	return books, nil
}

// Service

type BookService struct {
	userStore *UserStore
	bookStore *BookStore
}

func NewBookService(userStore *UserStore, bookStore *BookStore) *BookService {
	return &BookService{userStore: userStore, bookStore: bookStore}
}

func (b *BookService) SetRouter(r *mux.Router) {
	r.HandleFunc("/users/{userID}/books", b.booksHandler).Methods("GET", "POST")
	r.HandleFunc("/users/{userID}/books/{bookID}", b.bookHandler).Methods("GET", "DELETE")
}

func (b *BookService) booksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()

	switch r.Method {
	case "GET":
		user, err := b.userStore.Get(ctx, UserID(vars["userID"]))
		if err != nil {
			fmt.Fprintf(w, "get user error: %#v", err)
			return
		}
		fmt.Fprintf(w, "User: %#v\n", user)

		books, err := b.bookStore.List(ctx, user, "")
		if err != nil {
			fmt.Fprintf(w, "get book list error: %#v", err)
			return
		}
		if len(books) == 0 {
			fmt.Fprint(w, "no books!")
		}

		for i, book := range books {
			fmt.Fprintf(w, "#%d Book: %#v\n", i, book)
		}

		w.WriteHeader(http.StatusOK)
		return
	case "POST":
		bookName := r.FormValue("name")

		user, err := b.userStore.Get(ctx, UserID(vars["userID"]))
		if err != nil {
			fmt.Fprintf(w, "get user error: %#v", err)
			return
		}
		fmt.Fprintf(w, "User: %#v\n", user)

		book, err := NewBook(user, bookName)
		if err != nil {
			fmt.Fprintf(w, "create book error: %#v", err)
			return
		}

		err = b.bookStore.Put(ctx, user, book)
		if err != nil {
			fmt.Fprintf(w, "create book error: %#v", err)
			return
		}
		fmt.Fprintf(w, "Book: %#v\n", book)

		w.WriteHeader(http.StatusOK)
		return
	}
}

func (b *BookService) bookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()

	switch r.Method {
	case "GET":
		user, err := b.userStore.Get(ctx, UserID(vars["userID"]))
		if err != nil {
			fmt.Fprintf(w, "get user error: %#v", err)
			return
		}
		fmt.Fprintf(w, "User: %#v\n", user)

		book, err := b.bookStore.Get(ctx, user, BookID(vars["bookID"]))
		if err != nil {
			fmt.Fprintf(w, "get book error: %#v", err)
			return
		}
		fmt.Fprintf(w, "Book: %#v\n", book)

		w.WriteHeader(http.StatusOK)
		return
	case "DELETE":
		user, err := b.userStore.Get(ctx, UserID(vars["userID"]))
		if err != nil {
			fmt.Fprintf(w, "get user error: %#v", err)
			return
		}
		fmt.Fprintf(w, "User: %#v\n", user)

		book, err := b.bookStore.Get(ctx, user, BookID(vars["bookID"]))
		if err != nil {
			fmt.Fprintf(w, "get book error: %#v", err)
			return
		}

		err = b.bookStore.Delete(ctx, user, book.ID)
		if err != nil {
			fmt.Fprintf(w, "delete book error: %#v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}
