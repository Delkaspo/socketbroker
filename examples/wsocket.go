package main

import (
	"fmt"
	"github.com/SamuelRamond/socketbroker"
	sbplugin "github.com/SamuelRamond/socketbroker/clients"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const wcli = `
<html>
<body>
<div>Room 1: <span id="data_1">No data</span></div>
<div>Room 2: <span id="data_2">No data</span></div>

<script>
function ws(broker_uuid, data_id) {
    var socket = new WebSocket("ws://localhost:8080/ws/broker/"+broker_uuid);        

    socket.onopen = function() {
    	socket.send("Hello!");
    }

    socket.onmessage = function(message) {
		document.getElementById(data_id).innerHTML = message.data;
	}

	socket.onclose = function() {
		document.getElementById(data_id).innerHTML = "Connection to server lost...";
	}
}
ws("0xff42", "data_1");
ws("0xff21", "data_2");
</script>
</body>
</html>
`

func main() {
	b1 := socketbroker.New("Event 101", "0xff42")
	go b1.Run()

	b2 := socketbroker.New("Event 102", "0xff21")
	go b2.Run()

	// Register a client
	b1.Subscribe(&sbplugin.LogClient{})
	b2.Subscribe(&sbplugin.LogClient{})

	wsh := &sbplugin.SocketBrokerWsHandler{}
	wsh.Init()
	wsh.RegisterBroker(b1)
	wsh.RegisterBroker(b2)

	r := mux.NewRouter()

	r.HandleFunc("/ws/broker/{broker_uuid}", wsh.Handle).Methods("GET", "OPTIONS")

	r.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, wcli)

	}).Methods("GET", "OPTIONS")

	http.Handle("/", r)
	go func() {
		log.Fatal("ListenAndServe: ", http.ListenAndServe(":8080", nil))
	}()

	// Simple event Producer
	loop := 0
	for {
		<-time.After(time.Second)
		b1.Broadcast(fmt.Sprintf("hello world, loop: %d", loop))
		b2.Broadcast(fmt.Sprintf("EVXX, loop: %d", loop*2))
		loop++
	}
}
