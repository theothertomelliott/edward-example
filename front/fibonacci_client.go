package main

import (
	"net/rpc"
)

func getFibClient() (*rpc.Client, error) {
	var err error
	fc, err := rpc.DialHTTP("tcp", "127.0.0.1:8082")
	return fc, err
}

func currentFib() (int, error) {
	fibbClient, err := getFibClient()
	if err != nil {
		return 0, err
	}

	args := &struct{}{}
	var reply int
	err = fibbClient.Call("Fibonacci.Current", args, &reply)
	return reply, err
}

func nextFib() (int, error) {
	fibbClient, err := getFibClient()
	if err != nil {
		return 0, err
	}

	args := &struct{}{}
	var reply int
	err = fibbClient.Call("Fibonacci.Next", args, &reply)
	return reply, err
}
