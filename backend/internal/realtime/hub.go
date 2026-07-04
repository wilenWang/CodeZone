package realtime

import "sync"

type Hub struct {
	mu          sync.RWMutex
	connections map[int64]map[chan Event]struct{}
}

func NewHub() *Hub {
	return &Hub{connections: map[int64]map[chan Event]struct{}{}}
}

func (h *Hub) Register(userID int64, events chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[userID] == nil {
		h.connections[userID] = map[chan Event]struct{}{}
	}
	h.connections[userID][events] = struct{}{}
}

func (h *Hub) Unregister(userID int64, events chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userConnections := h.connections[userID]
	if userConnections == nil {
		return
	}
	delete(userConnections, events)
	if len(userConnections) == 0 {
		delete(h.connections, userID)
	}
}

func (h *Hub) SendToUser(userID int64, event Event) {
	h.mu.RLock()
	userConnections := make([]chan Event, 0, len(h.connections[userID]))
	for events := range h.connections[userID] {
		userConnections = append(userConnections, events)
	}
	h.mu.RUnlock()

	for _, events := range userConnections {
		events <- event
	}
}
