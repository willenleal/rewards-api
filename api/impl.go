package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Server struct {
	sync.RWMutex
	Receipts  map[string]Receipt
	Validator *validator.Validate
}

func NewServer() *Server {
	return &Server{
		Receipts:  make(map[string]Receipt),
		Validator: InitValidator(),
	}
}

type APIError struct {
	Description string `json:"description"`
}

type ReceiptProcessResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Server) PostReceiptsProcess(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()

	newReceipt := Receipt{}

	if err := json.NewDecoder(r.Body).Decode(&newReceipt); err != nil {
		slog.Error("err", "err", err)
		WriteJSON(w, http.StatusBadRequest, APIError{Description: "The receipt is invalid"})
		return
	}

	if err := s.Validator.Struct(newReceipt); err != nil {
		slog.Error("err", "err", err)
		WriteJSON(w, http.StatusBadRequest, APIError{Description: "The receipt is invalid"})
		return
	}

	uuid := uuid.New()

	s.Receipts[uuid.String()] = newReceipt

	resp := ReceiptProcessResponse{
		ID: uuid,
	}

	WriteJSON(w, http.StatusOK, resp)
}

type GetReceiptsIdPointsResponse struct {
	Points int64 `json:"points"`
}

func (s *Server) GetReceiptsIdPoints(w http.ResponseWriter, r *http.Request, id string) {

	receipt, ok := s.Receipts[id]

	if !ok {
		slog.Error("err", "err", "No receipt found for that id")
		WriteJSON(w, http.StatusNotFound, APIError{Description: "No receipt found for that id"})
		return
	}

	points, err := PointsCalculator(receipt)

	if err != nil {
		slog.Error("err", "err", err)
		WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	resp := GetReceiptsIdPointsResponse{
		Points: points,
	}

	WriteJSON(w, http.StatusOK, resp)
}

func WriteJSON(w http.ResponseWriter, status int, resp any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(resp)
}
