package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/realPointer/banknoteExchange/docs" // Swagger
	"github.com/realPointer/banknoteExchange/internal/service"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(l *zerolog.Logger, services *service.Services) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Minute))

	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/", NewExchangeRouter(services.Exchange, l))
	})

	return router
}
