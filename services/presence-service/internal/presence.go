package internal

import (
	"sync"
	"time"
)

type Presence struct {
	mu      sync.Mutex
	clients map[string]time.Time
}

func NewPresence() *Presence {
	return &Presence{
		clients: make(map[string]time.Time),
	}
}

// UpdatePresence records that a user is active at the current time.
func (p *Presence) UpdatePresence(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients[userID] = time.Now()
}

// GetOnlineUsers returns a list of user IDs that have been active in the last N seconds.
func (p *Presence) GetOnlineUsers(timeoutSeconds int) []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	online := []string{}
	threshold := time.Now().Add(-time.Duration(timeoutSeconds) * time.Second)
	for userID, lastSeen := range p.clients {
		if lastSeen.After(threshold) {
			online = append(online, userID)
		}
	}
	return online
}
