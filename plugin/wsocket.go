package plugin

import (
	"github.com/gorilla/websocket"
)

type WSocketClient struct {
	ws *websocket.Conn
	rcv chan interface{}
}

func (c *WSocketClient) Send(d interface{}) {
	c.rcv <- d
}

func (c *WSocketClient) Run() {
	for rcv := range c.rcv {
		err := c.ws.WriteJSON(rcv)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}