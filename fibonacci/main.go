package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	fibb := NewFibonacci()
	rpc.Register(fibb)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Fatal(http.Serve(l, nil))
}

func NewFibonacci() *Fibonacci {
	return &Fibonacci{
		previous: 0,
		current:  1,
	}
}

type Fibonacci struct {
	previous int
	current  int
}

func (f *Fibonacci) Current(req struct{}, res *int) error {
	log.Println("Returning current number")
	*res = f.current
	return nil
}

func (f *Fibonacci) Next(req struct{}, res *int) error {
	log.Println("Calculating next number")
	f.previous, f.current = f.current, f.previous+f.current
	*res = f.current
	return nil
}
