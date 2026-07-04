package realtime

import (
	"context"
	"sync"
)

const maxSubscriptionBacklog = 256

type Hub struct {
	mu          sync.RWMutex
	connections map[int64]map[*Subscription]struct{}
}

type Subscription struct {
	mu     sync.Mutex
	queue  []Event
	notify chan struct{}
	closed bool
}

func NewHub() *Hub {
	return &Hub{connections: map[int64]map[*Subscription]struct{}{}}
}

func (h *Hub) Register(userID int64) *Subscription {
	subscription := &Subscription{notify: make(chan struct{}, 1)}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[userID] == nil {
		h.connections[userID] = map[*Subscription]struct{}{}
	}
	h.connections[userID][subscription] = struct{}{}
	return subscription
}

func (h *Hub) Unregister(userID int64, subscription *Subscription) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userConnections := h.connections[userID]
	if userConnections == nil {
		return
	}
	delete(userConnections, subscription)
	subscription.close()
	if len(userConnections) == 0 {
		delete(h.connections, userID)
	}
}

func (h *Hub) SendToUser(userID int64, event Event) {
	h.mu.RLock()
	userConnections := make([]*Subscription, 0, len(h.connections[userID]))
	for subscription := range h.connections[userID] {
		userConnections = append(userConnections, subscription)
	}
	h.mu.RUnlock()

	for _, subscription := range userConnections {
		if !subscription.enqueue(event) {
			h.Unregister(userID, subscription)
		}
	}
}

func (s *Subscription) Next(ctx context.Context) (Event, bool) {
	for {
		s.mu.Lock()
		if s.closed {
			s.mu.Unlock()
			return Event{}, false
		}
		if len(s.queue) > 0 {
			event := s.queue[0]
			copy(s.queue, s.queue[1:])
			s.queue = s.queue[:len(s.queue)-1]
			s.mu.Unlock()
			return event, true
		}
		notify := s.notify
		s.mu.Unlock()

		select {
		case <-notify:
		case <-ctx.Done():
			return Event{}, false
		}
	}
}

func (s *Subscription) enqueue(event Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return false
	}
	if len(s.queue) >= maxSubscriptionBacklog {
		return false
	}
	s.queue = append(s.queue, event)
	select {
	case s.notify <- struct{}{}:
	default:
	}
	return true
}

func (s *Subscription) close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}
	s.closed = true
	s.queue = nil
	select {
	case s.notify <- struct{}{}:
	default:
	}
}
