package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "docker"
	DB_PASSWORD = "docker"
	DB_NAME     = "docker"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	msgs, err := NewMessages()
	if err != nil {
		log.Fatal(err)
	}
	rpc.Register(msgs)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Fatal(http.Serve(l, nil))
}

func NewMessages() (*Messages, error) {
	log.Println("Initializing database connection")
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	return &Messages{
		db: db,
	}, nil
}

type Messages struct {
	db *sql.DB
}

func (f *Messages) Get(req struct{}, res *[]string) error {
	log.Println("Getting messages")
	var msgs []string
	rows, err := f.db.Query("SELECT * FROM Message")
	if err != nil {
		return err
	}
	for rows.Next() {
		var msg string
		err = rows.Scan(&msg)
		if err != nil {
			return err
		}
		msgs = append(msgs, msg)
	}
	log.Println("Messages found:", len(msgs))
	*res = msgs
	return nil
}

func (f *Messages) Post(message string, res *struct{}) error {
	log.Println("Message:", message)
	*res = struct{}{}
	_, err := f.db.Exec("INSERT INTO Message(Message) VALUES($1);", message)
	if err != nil {
		return err
	}
	return nil
}

func (f *Messages) Clear(_ struct{}, _ *struct{}) error {
	log.Println("Clearing all messages")
	_, err := f.db.Exec("TRUNCATE TABLE Message;")
	return err
}
