package main

import (
	"net/rpc"
)

func getMessagesClient() (*rpc.Client, error) {
	var err error
	mc, err := rpc.DialHTTP("tcp", "127.0.0.1:8083")
	return mc, err
}

func getMessages() ([]string, error) {
	MessagesClient, err := getMessagesClient()
	if err != nil {
		return nil, err
	}
	var reply []string
	err = MessagesClient.Call("Messages.Get", struct{}{}, &reply)
	return reply, err
}

func postMessage(message string) error {
	MessagesClient, err := getMessagesClient()
	if err != nil {
		return err
	}
	err = MessagesClient.Call("Messages.Post", message, &struct{}{})
	return err
}

func clearMessages() error {
	MessagesClient, err := getMessagesClient()
	if err != nil {
		return err
	}
	return MessagesClient.Call("Messages.Clear", struct{}{}, &struct{}{})
}
