package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SlackClient struct {
	Token    string
	Channel  string
	Username string
	Icon     string
}

func (s *SlackClient) Send(d interface{}) {
	// @Todo: Json prettyfier .(type) map[..]
	payload, err := json.Marshal(map[string]string{
		"text":     fmt.Sprintf("```%v```", d),
		"channel":  s.Channel,
		"username": s.Username,
		"icon_url": s.Icon,
	})
	if err != nil {
		log.Println("SlackClient:", err)
		return
	}
	resp, err := http.Post(
		fmt.Sprintf(
			"https://hooks.slack.com/services/%s",
			s.Token,
		),
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		log.Println("SlackClient:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("SlackClient:", "There was an error warning slack...")
	}
}
