package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ticket-system/middleware"
	"ticket-system/models"
	"ticket-system/store"
	"ticket-system/utils"
)

type TicketHandler struct {
	Store *store.Store
}

type createTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

// CreateTicket creates a new ticket for the authenticated user.
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var req createTicketRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON",
		})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" || req.Description == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title and description are required",
		})
		return
	}

	ticket := &models.Ticket{
		ID:          utils.NewID(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
	}

	if ticket.ID == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to generate ticket ID",
		})
		return
	}

	h.Store.AddTicket(ticket)

	writeJSON(w, http.StatusCreated, ticket)
}

// ListTickets returns only tickets owned by the authenticated user.
func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	tickets := h.Store.GetTicketsByUser(userID)

	writeJSON(w, http.StatusOK, tickets)
}

// GetTicket returns a single ticket only if it belongs to the authenticated user.
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	ticketID := r.PathValue("id")
	if ticketID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "ticket ID is required",
		})
		return
	}

	ticket, exists := h.Store.GetTicket(ticketID)
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "ticket not found",
		})
		return
	}

	if ticket.UserID != userID {
		writeJSON(w, http.StatusForbidden, map[string]string{
			"error": "you do not have access to this ticket",
		})
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

// UpdateTicketStatus updates a ticket while enforcing the required status flow.
func (h *TicketHandler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	ticketID := r.PathValue("id")
	if ticketID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "ticket ID is required",
		})
		return
	}

	ticket, exists := h.Store.GetTicket(ticketID)
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "ticket not found",
		})
		return
	}

	if ticket.UserID != userID {
		writeJSON(w, http.StatusForbidden, map[string]string{
			"error": "you do not have access to this ticket",
		})
		return
	}

	var req updateStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON",
		})
		return
	}

	req.Status = strings.TrimSpace(req.Status)

	if req.Status != "open" &&
		req.Status != "in_progress" &&
		req.Status != "closed" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid status",
		})
		return
	}

	switch ticket.Status {
	case "open":
		if req.Status != "in_progress" {
			writeJSON(w, http.StatusConflict, map[string]string{
				"error": "open tickets can only move to in_progress",
			})
			return
		}

	case "in_progress":
		if req.Status != "closed" {
			writeJSON(w, http.StatusConflict, map[string]string{
				"error": "in_progress tickets can only move to closed",
			})
			return
		}

	case "closed":
		writeJSON(w, http.StatusConflict, map[string]string{
			"error": "closed tickets cannot be reopened or changed",
		})
		return
	}

	if err := h.Store.UpdateTicketStatus(ticketID, req.Status); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "ticket not found",
		})
		return
	}

	updatedTicket, _ := h.Store.GetTicket(ticketID)

	writeJSON(w, http.StatusOK, updatedTicket)
}
