package http

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HealthService interface {
	CheckDbHealth(ctx context.Context) error
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.HealthService.CheckDbHealth(r.Context()); err != nil {
		log.Info(fmt.Sprintf("unsuccessful request, status code: %d", http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		h.encodeJsonResponse(&w, responseEntity{Message: "Service Unhealthy"})
		return
	}
	log.Info(fmt.Sprintf("successfully handled request, status code: %d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	h.encodeJsonResponse(&w, responseEntity{Message: "Service alive. Database connection is good."})
	return
}
