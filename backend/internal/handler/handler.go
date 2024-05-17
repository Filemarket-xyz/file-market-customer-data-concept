package handler

import (
	"net/http"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/service"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
)

type Handler interface {
	Init() http.Handler
}

type handler struct {
	cfg               *config.HandlerConfig
	service           service.Service
	validationFormats strfmt.Registry
	logging           logger.Logger
}

func NewHandler(cfg *config.HandlerConfig, service service.Service, logging logger.Logger) Handler {
	return &handler{
		cfg:               cfg,
		service:           service,
		validationFormats: strfmt.NewFormats(),
		logging:           logging,
	}
}

func (h *handler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *handler) Init() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/config", h.GetConfig)

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Handle("/refresh", h.CookieRefreshAuthMiddleware((http.HandlerFunc(h.RefreshAuth))))
	authRouter.Handle("/logout", h.CookieAuthMiddleware((http.HandlerFunc(h.Logout))))
	authRouter.Handle("/full_logout", h.CookieAuthMiddleware((http.HandlerFunc(h.FullLogout))))
	authRouter.Handle("/try", h.CookieAuthMiddleware((http.HandlerFunc(h.TryAuth))))
	authRouter.HandleFunc("/message", h.AuthMessage)
	authRouter.HandleFunc("/by_signature", h.AuthByMessage)

	clientRouter := router.PathPrefix("/client").Subrouter()
	clientRouter.Handle("/dataset/download", h.CookieAuthMiddleware(h.permissionMiddleware(domain.RoleClient)((http.HandlerFunc(h.DownloadDataset)))))
	// clientRouter.HandleFunc("/dataset/download", h.DownloadDataset)
	clientRouter.Handle("/dataset/get", h.CookieAuthMiddleware(h.permissionMiddleware(domain.RoleClient)((http.HandlerFunc(h.GetDataset)))))
	clientRouter.Handle("/dataset/agreement", h.CookieAuthMiddleware(h.permissionMiddleware(domain.RoleClient)((http.HandlerFunc(h.AgreementDataset)))))

	router.Use(h.corsMiddleware)
	return router
}
