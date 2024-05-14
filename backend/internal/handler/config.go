package handler

import (
	"context"
	"net/http"
)

func (h *handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	res := h.service.GetConfig(ctx)

	if err := writeResponse(w, r, http.StatusOK, res); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}
