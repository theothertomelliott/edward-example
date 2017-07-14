package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.SetOutput(os.Stdout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fileServer = http.FileServer(http.Dir("public/"))

	http.HandleFunc("/messages/clear", ClearMessages)
	http.HandleFunc("/messages", GetMessages)
	http.HandleFunc("/fibonacci", GetFib)
	http.HandleFunc("/", Index)

	fmt.Println("Starting to listen on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var fileServer http.Handler

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Println("Serving file:", r.URL.Path)
		fileServer.ServeHTTP(w, r)
		return
	}
	log.Println("Serving Index")

	var err error
	t := template.New("index.html")
	t, err = t.ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var lastTouch string
	if r.Method == "POST" {
		_, err = touch()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}

	lastTouch, err = last()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"last": lastTouch,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetFib(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Fibonacci")

	var err error
	t := template.New("fibonacci.html")
	t, err = t.ParseFiles("views/fibonacci.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cf int
	if r.Method == "POST" {
		cf, err = nextFib()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}
	cf, err = currentFib()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"fibonacci": cf,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func ClearMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Message Clear")

	err := clearMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/messages", http.StatusFound)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Messages")

	var err error
	t := template.New("messages.html")
	t, err = t.ParseFiles("views/messages.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		msg := strings.Join(r.Form["message"], "")
		err = postMessage(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}
	var msgs []string
	msgs, err = getMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"messages": msgs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
