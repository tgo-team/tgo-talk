package tgo

import (
	"sync"
)

// -------------- clientManager -----------------------

type connManager struct {
	conns          map[uint64]Conn
	connLock       sync.RWMutex
	clientIDSequence int64
}

func newConnManager() *connManager {

	return &connManager{
		conns: make(map[uint64]Conn),
	}
}

func (cm *connManager) AddConn(connID uint64,conn Conn) uint64 {
	cm.connLock.Lock()
	cm.conns[connID] = conn
	cm.connLock.Unlock()
	return connID
}

func (cm *connManager) RemoveConn(connID uint64) {
	cm.connLock.Lock()
	_, ok := cm.conns[connID]
	if !ok {
		cm.connLock.Unlock()
		return
	}
	delete(cm.conns, connID)
	cm.connLock.Unlock()

}

func (cm *connManager) GetConn(connID uint64) Conn {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	return cm.conns[connID]
}
