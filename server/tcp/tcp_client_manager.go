package tcp

import (
	"github.com/tgo-team/tgo-talk/tgo"
	"sync"
)

// -------------- clientManager -----------------------

type connManager struct {
	conns          map[uint64]tgo.Conn
	connLock       sync.RWMutex
}

func newConnManager() *connManager {

	return &connManager{
		conns: make(map[uint64]tgo.Conn),
	}
}

func (cm *connManager) addConn(connID uint64,client tgo.Conn) uint64 {
	cm.connLock.Lock()
	cm.conns[connID] = client
	cm.connLock.Unlock()
	return connID
}

func (cm *connManager) removeConn(connID uint64) {
	cm.connLock.Lock()
	_, ok := cm.conns[connID]
	if !ok {
		cm.connLock.Unlock()
		return
	}
	delete(cm.conns, connID)
	cm.connLock.Unlock()

}

func (cm *connManager) getConn(clientId uint64) tgo.Conn {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	return cm.conns[clientId]
}
