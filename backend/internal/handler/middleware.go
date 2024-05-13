package handler

import (
	"context"
	"errors"

	"net/http"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/jwtoken"
)

type CtxKey string

const (
	CtxKeyUser = CtxKey("user")
	TokenStart = "Bearer "
)

func (h *handler) CookieAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cookie *http.Cookie
			err    error
		)
		cookie, err = r.Cookie(domain.AccessTokenCookieName)
		if err != nil || cookie == nil {
			h.logging.Info("cookie not found in ws connect")
			r.Body.Close()
			h.makeErrorResponse(w, r, errors.New("missing credentials"), code401)
			return
		}
		token := cookie.Value
		user, err := h.service.GetUserByJWToken(r.Context(), jwtoken.PurposeAccess, token)
		if err != nil {
			r.Body.Close()
			h.makeErrorResponse(w, r, err, code500)
			return
		}
		ctx := context.WithValue(r.Context(), CtxKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) CookieRefreshAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cookie *http.Cookie
			err    error
		)
		cookie, err = r.Cookie(domain.RefreshTokenCookieName)
		if err != nil || cookie == nil {
			h.logging.Info("cookie not found in ws connect")
			r.Body.Close()
			h.makeErrorResponse(w, r, errors.New("missing credentials"), code401)
			return
		}
		token := cookie.Value
		user, err := h.service.GetUserByJWToken(r.Context(), jwtoken.PurposeRefresh, token)
		if err != nil {
			r.Body.Close()
			h.makeErrorResponse(w, r, err, code500)
			return
		}
		ctx := context.WithValue(r.Context(), CtxKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) permissionMiddleware(roles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acc, ok := r.Context().Value(CtxKeyUser).(*domain.UserWithTokenNumber)
			if !ok {
				err := newHandlerError(
					code500,
					errors.New("permissionMiddleware: error: invalid user data in context: failed to convert to type *domain.UserWithTokenNumber"),
					"internal error",
					"invalid user data in context",
				)
				h.makeErrorResponse(w, r, err, code500)
				return
			}
			for _, role := range roles {
				if role == acc.Role {
					next.ServeHTTP(w, r)
					return
				}
			}
			h.makeErrorResponse(w, r, errors.New("permission denied"), http.StatusForbidden)
		})
	}
}
