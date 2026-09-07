package store

import (
	"errors"
	"sync"

	"ticket-system/models"
)

// Store keeps application data in memory.
// RWMutex makes access safe when multiple requests arrive at the same time.
type Store struct {
	mu      sync.RWMutex
	users   map[string]*models.User
	tickets map[string]*models.Ticket
}

// NewStore creates an empty in-memory store.
func NewStore() *Store {
	return &Store{
		users:   make(map[string]*models.User),
		tickets: make(map[string]*models.Ticket),
	}
}

// AddUser stores a new user.
func (s *Store) AddUser(user *models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.ID] = user
}

// GetUserByEmail finds a user by email.
func (s *Store) GetUserByEmail(email string) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, true
		}
	}

	return nil, false
}

// AddTicket stores a new ticket.
func (s *Store) AddTicket(ticket *models.Ticket) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tickets[ticket.ID] = ticket
}

// GetTicket returns a ticket by ID.
func (s *Store) GetTicket(id string) (*models.Ticket, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ticket, ok := s.tickets[id]
	return ticket, ok
}

// GetTicketsByUser returns only tickets belonging to a specific user.
func (s *Store) GetTicketsByUser(userID string) []*models.Ticket {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tickets := make([]*models.Ticket, 0)

	for _, ticket := range s.tickets {
		if ticket.UserID == userID {
			tickets = append(tickets, ticket)
		}
	}

	return tickets
}

// UpdateTicketStatus changes the status of a ticket.
func (s *Store) UpdateTicketStatus(id string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, ok := s.tickets[id]
	if !ok {
		return errors.New("ticket not found")
	}

	ticket.Status = status
	return nil
}
