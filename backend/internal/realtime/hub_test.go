package realtime

import "testing"

func TestHubStoresConnectionsByUser(t *testing.T) {
	hub := NewHub()
	events := make(chan Event, 1)

	hub.Register(42, events)
	hub.SendToUser(42, Event{Type: "message.created", Payload: map[string]any{"id": int64(7)}})

	got := <-events
	if got.Type != "message.created" {
		t.Fatalf("got event type %q want message.created", got.Type)
	}

	hub.Unregister(42, events)
	hub.SendToUser(42, Event{Type: "message.created"})

	select {
	case event := <-events:
		t.Fatalf("got event after unregister: %#v", event)
	default:
	}
}
