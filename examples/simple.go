package main 

import (
	"time"
	sbplugin "github.com/SamuelRamond/socketbroker/plugin"
	"github.com/SamuelRamond/socketbroker"
)

func main() {
	b := socketbroker.New("Event 101")
	go b.Run()

	// Register a client
	b.Subscribe(&sbplugin.LogClient{})
	b.Subscribe(&sbplugin.HttpGetClient{
		Url: "https://golang.org",
	})

	// Simple event Producer
	for {
		<-time.After(time.Second)
		b.Broadcast("hello world")
		
		<-time.After(time.Second)
		b.Broadcast(map[string]interface{}{
			"ev": "login",
			"data": "sam@yolo connected",	
		})
	}
}
