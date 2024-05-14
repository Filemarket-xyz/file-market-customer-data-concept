package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
)

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(CtxKeyUser).(*domain.UserWithTokenNumber)
	if !ok {
		h.makeErrorResponse(w, r, UserMissingInCtxErr, code500)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	if err := h.service.Logout(ctx, user); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	setCookies(w, &domain.PairOfTokens{
		AccessToken:  &domain.Token{},
		RefreshToken: &domain.Token{},
	})
	result := true
	if err := writeResponse(w, r, http.StatusOK, &models.SuccessResponse{Success: &result}); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}

func (h *handler) FullLogout(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(CtxKeyUser).(*domain.UserWithTokenNumber)
	if !ok {
		h.makeErrorResponse(w, r, UserMissingInCtxErr, code500)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	if err := h.service.FullLogout(ctx, user); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	setCookies(w, &domain.PairOfTokens{
		AccessToken:  &domain.Token{},
		RefreshToken: &domain.Token{},
	})
	result := true
	if err := writeResponse(w, r, http.StatusOK, &models.SuccessResponse{Success: &result}); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}

func (h *handler) AuthMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req models.AuthMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	if err := req.Validate(h.validationFormats); err != nil {
		h.makeErrorResponse(w, r, makeValidationError("handleAuthMessage", err), code400)
		return
	}
	if err := addressValidation(*req.Address); err != nil {
		h.makeErrorResponse(w, r, makeValidationError("handleAuthMessage", err), code400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	res, err := h.service.GetAuthMessage(ctx, &req)
	if err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	if err := writeResponse(w, r, http.StatusOK, res); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}

func (h *handler) AuthByMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req models.AuthBySignatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	if err := req.Validate(h.validationFormats); err != nil {
		h.makeErrorResponse(w, r, makeValidationError("AuthByMessage", err), code400)
		return
	}
	if err := addressValidation(*req.Address); err != nil {
		h.makeErrorResponse(w, r, makeValidationError("AuthByMessage", err), code400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	res, tokens, err := h.service.AuthByMessage(ctx, &req)
	if err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	setCookies(w, tokens)
	if err := writeResponse(w, r, http.StatusOK, res); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}

func (h *handler) RefreshAuth(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	user, ok := ctx.Value(CtxKeyUser).(*domain.UserWithTokenNumber)
	if !ok {
		h.makeErrorResponse(w, r, errors.New("RefreshAuth"), code500)
		return
	}

	res, tokens, err := h.service.RefreshJWTokens(ctx, user)
	if err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
	setCookies(w, tokens)
	if err := writeResponse(w, r, http.StatusOK, res); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}

func (h *handler) TryAuth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()
	user, ok := ctx.Value(CtxKeyUser).(*domain.UserWithTokenNumber)
	if !ok {
		h.makeErrorResponse(w, r, errors.New("TryAuth"), code500)
		return
	}

	res, err := h.service.AuthByToken(ctx, user)
	if err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}

	if err := writeResponse(w, r, http.StatusOK, res); err != nil {
		h.makeErrorResponse(w, r, err, code500)
		return
	}
}

func setCookies(w http.ResponseWriter, tokens *domain.PairOfTokens) {
	http.SetCookie(w, &http.Cookie{
		Name:     domain.AccessTokenCookieName,
		Value:    tokens.AccessToken.Token,
		Path:     "/",
		Expires:  tokens.AccessToken.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     domain.RefreshTokenCookieName,
		Value:    tokens.RefreshToken.Token,
		Path:     "/api/auth/refresh",
		Expires:  tokens.RefreshToken.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
