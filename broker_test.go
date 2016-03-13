package socketbroker

import (
	"testing"
	"time"
)

type TestClient struct {
	LastCall interface{}
}

func (t *TestClient) Send(i interface{}) {
	t.LastCall = i
}

func TestBroker(t *testing.T) {
	s := New("test", "0xtest")

	go s.Run()

	clitest := &TestClient{}
	s.Subscribe(clitest)

	s.Broadcast(1)

	<-time.After(time.Millisecond)
	if clitest.LastCall.(int) != 1 {
		t.Error("Broadcast failed")
	}

	s.Broadcast("hello")
	<-time.After(time.Millisecond)
	if clitest.LastCall.(string) != "hello" {
		t.Error("Broadcast failed")
	}

	s.Unsubscribe(clitest)
	<-time.After(time.Millisecond)

	s.Broadcast("Ola")
	<-time.After(time.Millisecond)
	if clitest.LastCall.(string) != "hello" {
		t.Error("Broadcast should not have worked on this client")
	}
}
