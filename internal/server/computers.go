package server

import (
	"encoding/json"
	"github.com/ninja-way/pc-store/internal/model"
	"io"
	"net/http"
	"path"
	"strconv"
)

func (s Server) getComputers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	comps, err := s.db.GetComputers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(comps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(res)
}

func (s Server) addComputer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var newPC model.PC
	if err = json.Unmarshal(body, &newPC); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.db.AddComputer(newPC); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s Server) manageComputer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.RequestURI))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getComputer(w, r, id)
	case http.MethodPost:
		s.updateComputer(w, r, id)
	case http.MethodDelete:
		s.deleteComputer(w, r, id)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (s Server) getComputer(w http.ResponseWriter, _ *http.Request, id int) {
	pc, err := s.db.GetComputerByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(pc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(res)
}

func (s Server) updateComputer(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var newPC model.PC
	if err = json.Unmarshal(body, &newPC); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.db.UpdateComputer(id, newPC); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s Server) deleteComputer(w http.ResponseWriter, _ *http.Request, id int) {
	if err := s.db.DeleteComputer(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}
