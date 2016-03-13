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

const wcli = `<html>
<body>
<h3>Socket Broker Demo</h3>
<h5>Room 1: (Producer is server hearthbeat)</h5>
<p id="data_1">No data</p>
<h5>Room 2: (Producer is js client)</h5>
<p id="data_2">No data</p>

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
	return socket;
}
ws("0xff42", "data_1");
var r2 = ws("0xff21", "data_2");
var i = 0;
setInterval(function () {r2.send("Hello world, loop:"+i);i++;}, 1500);
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

	// Create a WShandler
	wsh := &sbplugin.SocketBrokerWsHandler{}
	wsh.Init(sbplugin.SbWsHReadWrite)
	
	// Register your brokers
	wsh.RegisterBroker(b1)
	wsh.RegisterBroker(b2)

	r := mux.NewRouter()

	// Register the handler in the router
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
		b1.Broadcast(fmt.Sprintf("Hello world, loop: %d", loop))
		loop++
	}
}
