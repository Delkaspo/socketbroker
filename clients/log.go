package clients

import (
	"log"
)

type LogClient struct{}

func (l *LogClient) Send(v interface{}) {
	log.Println(v)
}
