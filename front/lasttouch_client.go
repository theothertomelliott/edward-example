package main

import (
	"log"
	"net/rpc"
)

var lastTouchClient *rpc.Client

func init() {
	var err error
	lastTouchClient, err = rpc.DialHTTP("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatal("dialing:", err)
	}
}

func last() (string, error) {
	args := &struct{}{}
	var reply string
	err := lastTouchClient.Call("LastTouch.Last", args, &reply)
	return reply, err
}

func touch() (string, error) {
	args := &struct{}{}
	var reply string
	err := lastTouchClient.Call("LastTouch.Touch", args, &reply)
	return reply, err
}
