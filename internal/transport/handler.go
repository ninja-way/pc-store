package transport

import (
	"context"
	"encoding/json"
	"github.com/ninja-way/pc-store/internal/models"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
)

type ComputersStore interface {
	GetComputers(context.Context) ([]models.PC, error)
	GetComputerByID(context.Context, int) (models.PC, error)
	AddComputer(context.Context, models.PC) error
	UpdateComputer(context.Context, int, models.PC) error
	DeleteComputer(context.Context, int) error
}

type Handler struct {
	service ComputersStore
}

func NewHandler(service ComputersStore) *Handler {
	return &Handler{
		service: service,
	}
}

// Get method which returns all pc from database
func (h *Handler) getComputers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	comps, err := h.service.GetComputers(context.TODO())
	if err != nil {
		log.Println("GetComputers error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(comps)
	if err != nil {
		log.Println("GetComputers error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(res)
}

// Put method which add new pc from request body to database
func (h *Handler) addComputer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newPC models.PC
	if err = json.Unmarshal(body, &newPC); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.service.AddComputer(context.TODO(), newPC); err != nil {
		log.Println("AddComputer error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Manage check which method came to the same ID endpoint and parse ID
func (h *Handler) manageComputer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.RequestURI))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getComputer(w, r, id)
	case http.MethodPost:
		h.updateComputer(w, r, id)
	case http.MethodDelete:
		h.deleteComputer(w, r, id)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// Get method which return pc from database by id
func (h *Handler) getComputer(w http.ResponseWriter, _ *http.Request, id int) {
	pc, err := h.service.GetComputerByID(context.TODO(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(pc)
	if err != nil {
		log.Println("GetComputer error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(res)
}

// Post method which update existing pc in database by id and new data from request body
func (h *Handler) updateComputer(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newPC models.PC
	if err = json.Unmarshal(body, &newPC); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.service.UpdateComputer(context.TODO(), id, newPC); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Delete method which delete pc from database by id
func (h *Handler) deleteComputer(w http.ResponseWriter, _ *http.Request, id int) {
	if err := h.service.DeleteComputer(context.TODO(), id); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}
