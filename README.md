Socket Broker
-------------

Broadcast events to connected parties.

### Summary:

```
									+------> Websocket
				  					|
[data producer]---->SocketBroker----+------> HTTPPost
									|
									+------> Slack
									|
									+------> Cassandra
```

### Usage:

See examples folder for real examples.

- Simple handler
- Websocket multi broker handler

Simple:
```
package main 

import (
	"time"
	"github.com/SamuelRamond/socketbroker/clients"
	"github.com/SamuelRamond/socketbroker"
)

func main() {
	b := socketbroker.New("Event 101")
	go b.Run()

	// Register clients
	b.Subscribe(&clients.LogClient{})
	b.Subscribe(&clients.HttpGetClient{
		Url: "http://url.to.your.hook",
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
```
output:

```
2016/03/11 22:02:45 New subscriber ! (live 1)
2016/03/11 22:37:54 hello world
2016/03/11 22:37:55 map[data:sam@yolo connected ev:login]
```

### Broker clients

- [x] Log
- [x] HttpPostHook
- [x] HttpGetHook
- [x] Slack
- [x] Websocket
- [ ] Protobuf
- [ ] Rmq
- [ ] Cassandra

## Producers

- [ ] FS
- [ ] HTTP
- [ ] Slack
- [ ] Http poller

```
Author: @erazor42
```


