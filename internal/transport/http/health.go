package http

import (
	"context"
	"net/http"
)

type HealthService interface {
	CheckDbHealth(ctx context.Context) error
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.HealthService.CheckDbHealth(r.Context()); err != nil {
		h.SendServerErrorResponse(w, r, err)
		return
	}
	message := Response{Message: "Service alive. Database connection is good."}
	h.SendOkResponse(w, r, message)
	return
}
