package passenger

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hovanja2011/move-together/internal/apperror"
	"github.com/hovanja2011/move-together/internal/handlers"
	"github.com/hovanja2011/move-together/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

const (
	passengersURL = "/passengers"
	passengerURL  = "/passenger/:uuid"
)

type handler struct {
	repository Repository
	logger     *logging.Logger
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, passengersURL, apperror.Middleware(h.GetList))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		w.WriteHeader(400)
		return err
	}
	allBytes, err := json.Marshal(all)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)
	return nil
}
