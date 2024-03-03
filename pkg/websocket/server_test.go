package websocket

import (
	"testing"
)

func TestManager(t *testing.T) {
	manager := &Manager{
		Register:    make(chan *Client, 128),
		UnRegister:  make(chan *Client, 128),
		Message:     make(chan *MessageData, 128),
		clientCount: 0,
	}

	client := &Client{
		Id: "testClient",
	}

	// Test RegisterClient
	manager.RegisterClient(client)
	registeredClient, ok := manager.Group.Load(client.Id)
	if !ok || registeredClient.(*Client).Id != client.Id {
		t.Errorf("RegisterClient failed, client not found in manager's group")
	}

	// Test UnRegisterClient
	manager.UnRegisterClient(client)
	_, ok = manager.Group.Load(client.Id)
	if ok {
		t.Errorf("UnRegisterClient failed, client still found in manager's group")
	}
}
