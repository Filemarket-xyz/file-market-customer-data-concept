package handler

import (
	"context"
	"encoding/csv"
	"errors"
	"net/http"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
)

func (h *handler) DownloadDataset(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	user, ok := ctx.Value(CtxKeyUser).(*domain.UserWithTokenNumber)
	if !ok {
		h.makeErrorResponse(w, r, errors.New("DownloadDataset"), code500)
		return
	}

	datasets, err := h.service.GetUserDataset(ctx, user)
	if err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}

	// Set our headers so browser will download the file
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=dataset.csv")
	// Create a CSV writer using our HTTP response writer as our io.Writer
	wr := csv.NewWriter(w)
	// Write all items and deal with errors
	if err := wr.WriteAll(datasets); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}
