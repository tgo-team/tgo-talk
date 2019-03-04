package server

import (
	"sync"
)

// -------------- clientManager -----------------------

type clientManager struct {
	clients          map[int64]*Client
	clientLock       sync.RWMutex
	clientIDSequence int64
}

func newClientManager() *clientManager {

	return &clientManager{
		clients: make(map[int64]*Client),
	}
}

func (cm *clientManager) addClient(clientId int64,client *Client) int64 {
	client.id = clientId
	cm.clientLock.Lock()
	cm.clients[client.id] = client
	cm.clientLock.Unlock()
	return clientId
}

func (cm *clientManager) removeClient(clientId int64) {
	cm.clientLock.Lock()
	_, ok := cm.clients[clientId]
	if !ok {
		cm.clientLock.Unlock()
		return
	}
	delete(cm.clients, clientId)
	cm.clientLock.Unlock()

}

func (cm *clientManager) getClient(clientId int64) (cli *Client) {
	cm.clientLock.Lock()
	cli = cm.clients[clientId]
	cm.clientLock.Unlock()
	return
}
