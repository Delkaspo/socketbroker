package clients

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type HttpPostClient struct {
	Url string
}

// Post the payload to the configured url
func (s *HttpPostClient) Send(d interface{}) {
	payload, err := json.Marshal(d)
	if err != nil {
		log.Println("HttpPostClient:", err)
		return
	}
	resp, err := http.Post(
		s.Url,
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		log.Println("HttpPostClient:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("HttpPostClient:", "HTTP::POST Error", resp.StatusCode)
	}
}

type HttpGetClient struct {
	Url string
}

// Post the payload to the configured url
func (s *HttpGetClient) Send(d interface{}) {
	resp, err := http.Get(
		s.Url,
	)
	if err != nil {
		log.Println("HttpGetClient:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("HttpGetClient:", "HTTP::GET Error", resp.StatusCode)
	}
}
