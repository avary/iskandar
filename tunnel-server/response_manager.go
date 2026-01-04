package main

import (
	"sync"

	"github.com/igneel64/iskandar/shared/protocol"
)

type MessageChannel = chan protocol.Message

type RequestManager interface {
	GetRequestChannel(requestId string) (MessageChannel, bool)
	RegisterRequest(requestId string) MessageChannel
	RemoveRequest(requestId string)
}

type InMemoryRequestManager struct {
	requestChannelMap map[string]MessageChannel
	mu                sync.RWMutex
}

func NewInMemoryRequestManager() *InMemoryRequestManager {
	return &InMemoryRequestManager{
		requestChannelMap: make(map[string]MessageChannel),
	}
}

func (i *InMemoryRequestManager) GetRequestChannel(requestId string) (MessageChannel, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	ch, ok := i.requestChannelMap[requestId]
	return ch, ok
}

func (i *InMemoryRequestManager) RegisterRequest(requestId string) MessageChannel {
	i.mu.Lock()
	defer i.mu.Unlock()
	requestChannel := make(MessageChannel)
	i.requestChannelMap[requestId] = requestChannel
	return requestChannel
}

func (i *InMemoryRequestManager) RemoveRequest(requestId string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if ch, ok := i.requestChannelMap[requestId]; ok {
		close(ch)
		delete(i.requestChannelMap, requestId)
	}
}
