package clients

import (
	"github.com/SamuelRamond/socketbroker"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

type WSocketClient struct {
	ws *websocket.Conn
	out chan interface{}
	bcast chan interface{}
}

func (c *WSocketClient) Send(d interface{}) {
	c.out <- d
}

func (c *WSocketClient) fin() {
	for {
		_, in, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		if c.bcast != nil {
			c.bcast <- string(in)
		}
	}
	c.ws.Close()
}

func (c *WSocketClient) fout() {
	for out := range c.out {
		err := c.ws.WriteJSON(out)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func (c *WSocketClient) Run() {
	go c.fin()
	c.fout()
}

const (
	SbWsHReadOnly = iota
	SbWsHReadWrite
)

type SocketBrokerWsHandler struct {
	upg *websocket.Upgrader
	brokers map[string]*socketbroker.SocketBroker
	mode int
}

func (c *SocketBrokerWsHandler) Init(mode int) {
	c.upg = &websocket.Upgrader{ReadBufferSize: 2048, WriteBufferSize: 2048}
	c.brokers = make(map[string]*socketbroker.SocketBroker)
	c.mode = mode
}

// Register a broker to a ws handler so any client can get broker event via 
// websocket
func (c *SocketBrokerWsHandler) RegisterBroker(s *socketbroker.SocketBroker) {
	log.Printf("SocketBrokerWsHandler: Registering broker: %s\n", s.UUID)
	c.brokers[s.UUID] = s
}

// Handle and dispatch websocket connection, this require to have an url matching:
// `/what/ever/{broker_uuid}`
func (c *SocketBrokerWsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Check if broker match uuid
	params := mux.Vars(r)
	buuid, ok := params["broker_uuid"]
	if !ok {
		log.Printf("Nobroker matching uuid: %s", buuid)
		return
	}

	// Upgrade websocket
	ws, err := c.upg.Upgrade(w, r, nil)
	if err != nil {
		// @Todo: Handle proper error
		log.Println(err)
		return
	}
	// Build client
	var chanin chan interface{} 
	if c.mode == SbWsHReadWrite {
		chanin = c.brokers[buuid].Bcast
	}
	cli := &WSocketClient{
		ws,
		make(chan interface{}),
		chanin,
	}

	c.brokers[buuid].Subscribe(cli)
	cli.Run()
	c.brokers[buuid].Unsubscribe(cli)
}
