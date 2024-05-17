package handler

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
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
	// user := &domain.UserWithTokenNumber{
	// 	Id: 4,
	// 	Role: domain.RoleClient,
	// }

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

func (h *handler) GetDataset(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	user, ok := ctx.Value(CtxKeyUser).(*domain.UserWithTokenNumber)
	if !ok {
		h.makeErrorResponse(w, r, errors.New("GetDataset"), code500)
		return
	}

	res, err := h.service.GetDataset(ctx, user.Id)
	if err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	if err := writeResponse(w, r, http.StatusOK, res); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}

func (h *handler) AgreementDataset(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	user, ok := ctx.Value(CtxKeyUser).(*domain.UserWithTokenNumber)
	if !ok {
		h.makeErrorResponse(w, r, errors.New("AgreementDataset"), code500)
		return
	}
	var req models.ClientAgreementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	if err := req.Validate(h.validationFormats); err != nil {
		h.makeErrorResponse(w, r, makeValidationError("AgreementDataset", err), code400)
		return
	}

	res, err := h.service.UpdateClientAgreement(ctx, user.Id, *req.Agreement)
	if err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	if err := writeResponse(w, r, http.StatusOK, res); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}
