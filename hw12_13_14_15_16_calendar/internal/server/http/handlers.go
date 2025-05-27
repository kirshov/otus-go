package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/domain"
)

type Handler struct {
	app Application
}
type jsonResponse map[string]interface{}

func NewHandlers(app Application) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, err := w.Write([]byte("Main page"))
	if err != nil {
		h.app.GetLogger().Error("Write failed")
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, h, http.MethodPost) {
		return
	}

	event, ok := parseRequest(w, r, h)
	if !ok {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.app.GetLogger().Error("Body.Close() failed")
		}
	}(r.Body)

	ctx := context.Background()

	if _, err := h.app.GetStorage().Add(ctx, event); err != nil {
		h.app.GetLogger().Error("event add failed: %w" + err.Error())
		sendResponse(w, h, http.StatusBadRequest, jsonResponse{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	sendResponse(w, h, http.StatusCreated, jsonResponse{"status": true})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, h, http.MethodPut) {
		return
	}

	event, ok := parseRequest(w, r, h)
	if !ok {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.app.GetLogger().Error("Body.Close() failed")
		}
	}(r.Body)

	ctx := context.Background()
	if _, err := h.app.GetStorage().GetByID(ctx, event.ID); err != nil {
		sendResponse(w, h, http.StatusNotFound, jsonResponse{})
		return
	}

	if err := h.app.GetStorage().Update(ctx, event); err != nil {
		h.app.GetLogger().Error("event update failed: %w" + err.Error())
		sendResponse(w, h, http.StatusBadRequest, jsonResponse{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	sendResponse(w, h, http.StatusOK, jsonResponse{"status": true})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, h, http.MethodDelete) {
		return
	}

	event, ok := parseRequest(w, r, h)
	if !ok {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.app.GetLogger().Error("Body.Close() failed")
		}
	}(r.Body)

	ctx := context.Background()
	if err := h.app.GetStorage().Remove(ctx, event.ID); err != nil {
		sendResponse(w, h, http.StatusNotFound, jsonResponse{})
		return
	}

	sendResponse(w, h, http.StatusOK, jsonResponse{"status": true})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, h, http.MethodGet) {
		return
	}

	daysStr := r.URL.Query().Get("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		sendResponse(w, h, http.StatusBadRequest, jsonResponse{
			"status": false,
			"error":  "days must be an integer",
		})
		return
	}

	ctx := context.Background()
	list, err := h.app.GetStorage().List(ctx, days)
	if err != nil {
		h.app.GetLogger().Error(err.Error())
		sendResponse(w, h, http.StatusInternalServerError, jsonResponse{})
		return
	}

	sendResponse(w, h, http.StatusOK, jsonResponse{"status": true, "items": list})
}

func checkMethod(w http.ResponseWriter, r *http.Request, h *Handler, m string) bool {
	if r.Method != m {
		sendResponse(w, h, http.StatusBadRequest, jsonResponse{
			"status": false,
			"error":  "support only POST method",
		})
		return false
	}

	return true
}

func parseRequest(w http.ResponseWriter, r *http.Request, h *Handler) (domain.Event, bool) {
	event := domain.Event{}

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		sendResponse(w, h, http.StatusBadRequest, jsonResponse{
			"status":      false,
			"error":       "Invalid request",
			"description": err.Error(),
		})
		return event, false
	}

	return event, true
}

func sendResponse(w http.ResponseWriter, h *Handler, s int, r jsonResponse) {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")

	if len(r) == 0 {
		return
	}

	message, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		h.app.GetLogger().Error("marshal failed: " + err.Error())
		return
	}

	if _, err := w.Write(message); err != nil {
		h.app.GetLogger().Error("http write failed: " + err.Error())
	}
}
