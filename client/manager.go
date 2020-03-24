package client

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// Manager is used to represent an overview of all currently
// connected clients
type Manager struct {
	clients []*Client
	sync.Mutex
}

// List returns all currently tracked clients
func (m *Manager) List() []*Client {
	return m.clients
}

// Client - return for given id or nil
// If UUID is 0, a new one will be generated
func (m *Manager) Client(id uint32) *Client {
	for _, c := range m.clients {
		if c.UUID == id {
			return c
		}
	}

	// If we're given a bad ID, don't create a new one
	// this can happen if we lag behind a client close command
	if id > 0 {
		return nil
	}

	newid := uuid.New()
	id = newid.ID()

	// In the case of client applications, we may end up with several active
	// clients connecting at a time. The special buffer "none" denotes that
	// the client is connected, and won't have to worry about unread count inaccuracy
	client := &Client{
		UUID:    id,
		current: "none",
	}

	// Lock around the append in case we're accessing in another thread
	m.Lock()
	m.clients = append(m.clients, client)
	m.Unlock()

	return client
}

// Remove a named grieth from the internal list
func (m *Manager) Remove(uuid uint32) error {
	for n, c := range m.clients {
		if c.UUID == uuid {
			// Lock around update incase we're accessing in another thread
			m.Lock()
			m.clients[n] = m.clients[len(m.clients)-1]
			m.clients = m.clients[:len(m.clients)-1]
			m.Unlock()

			return nil
		}
	}

	return errors.New("no client for UUID")
}
