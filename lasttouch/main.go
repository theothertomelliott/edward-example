package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	lt := NewLastTouch()
	rpc.Register(lt)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Fatal(http.Serve(l, nil))
}

func NewLastTouch() *LastTouch {
	return &LastTouch{}
}

type LastTouch struct {
	t *time.Time
}

func (f *LastTouch) Last(req struct{}, res *string) error {
	log.Println("Returning last touch time")
	if f.t == nil {
		*res = "Never"
	} else {
		*res = f.t.Format("2006-01-02 15:04:05")
	}
	return nil
}

func (f *LastTouch) Touch(req struct{}, res *string) error {
	t := time.Now()
	f.t = &t
	log.Println("Time is now:", t)
	return f.Last(req, res)
}
