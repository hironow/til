package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
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
		log.Printf("Book: %+v", book)
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
			log.Printf("#%d Book: %+v", i, book)
		}
	}

	return books, nil
}

// Service

type BookService struct {
	handlerFuncMap map[string]http.HandlerFunc
}

func NewBookService(userStore *UserStore, bookStore *BookStore) *BookService {
	m := make(map[string]http.HandlerFunc)

	// TODO: users/xxx/books/yyy
	m["/books/"] = func(w http.ResponseWriter, r *http.Request) {
		//ctx := r.Context()
	}

	return &BookService{handlerFuncMap: m}
}
