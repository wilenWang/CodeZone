package realtime

import (
	"context"
	"testing"
	"time"
)

func TestHubStoresConnectionsByUser(t *testing.T) {
	hub := NewHub()
	subscription := hub.Register(42)

	hub.SendToUser(42, Event{Type: "message.created", Payload: map[string]any{"id": int64(7)}})

	got, ok := subscription.Next(context.Background())
	if !ok {
		t.Fatal("subscription closed before receiving event")
	}
	if got.Type != "message.created" {
		t.Fatalf("got event type %q want message.created", got.Type)
	}

	hub.Unregister(42, subscription)
	hub.SendToUser(42, Event{Type: "message.created"})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if event, ok := subscription.Next(ctx); ok {
		t.Fatalf("got event after unregister: %#v", event)
	}
}

func TestSendToUserDoesNotBlockWhenSubscriptionHasPendingEvents(t *testing.T) {
	hub := NewHub()
	subscription := hub.Register(42)
	defer hub.Unregister(42, subscription)

	for i := 0; i < maxSubscriptionBacklog; i++ {
		hub.SendToUser(42, Event{Type: "message.created", Payload: i})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 0; i < maxSubscriptionBacklog; i++ {
		event, ok := subscription.Next(ctx)
		if !ok {
			t.Fatalf("subscription closed after %d events", i)
		}
		if event.Payload != i {
			t.Fatalf("got payload %#v want %d", event.Payload, i)
		}
	}
}

func TestSendToUserDisconnectsSubscriptionImmediatelyWhenBacklogIsFull(t *testing.T) {
	hub := NewHub()
	subscription := hub.Register(42)

	for i := 0; i < maxSubscriptionBacklog; i++ {
		hub.SendToUser(42, Event{Type: "message.created", Payload: i})
	}
	hub.SendToUser(42, Event{Type: "message.created", Payload: maxSubscriptionBacklog})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if event, ok := subscription.Next(ctx); ok {
		t.Fatalf("got event after backlog disconnect: %#v", event)
	}
}
