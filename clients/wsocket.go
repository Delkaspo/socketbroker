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
}

func (c *WSocketClient) Send(d interface{}) {
	c.out <- d
}

func (c *WSocketClient) fin() {
	for {
		n, in, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		log.Println(n, string(in))
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

type SocketBrokerWsHandler struct {
	upg *websocket.Upgrader
	brokers map[string]*socketbroker.SocketBroker
}

func (c *SocketBrokerWsHandler) Init() {
	c.upg = &websocket.Upgrader{ReadBufferSize: 2048, WriteBufferSize: 2048}
	c.brokers = make(map[string]*socketbroker.SocketBroker)
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
	cli := &WSocketClient{ws, make(chan interface{})}

	c.brokers[buuid].Subscribe(cli)
	cli.Run()
	c.brokers[buuid].Unsubscribe(cli)
}
