package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/realPointer/banknoteExchange/internal/service"
	"github.com/rs/zerolog"
)

type exchangeRoutes struct {
	exchangeService service.Exchange
	l               *zerolog.Logger
}

func NewExchangeRouter(exchangeService service.Exchange, l *zerolog.Logger) http.Handler {
	s := &exchangeRoutes{
		exchangeService: exchangeService,
		l:               l,
	}

	router := chi.NewRouter()
	router.Post("/exchange", s.exchangeBanknotes)

	return router
}

type ExchangeRequest struct {
	Amount    int   `json:"amount"    validate:"required,min=1"`
	Banknotes []int `json:"banknotes" validate:"required,min=1,unique,dive,min=1"`
}

type ExchangeResponse struct {
	Exchanges [][]int `json:"exchanges"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// @Summary Exchange banknotes
// @Description Exchange banknotes
// @Tags exchange
// @Accept json
// @Produce json
// @Param request body ExchangeRequest true "Exchange request"
// @Success 200 {object} ExchangeResponse
// @Failure 400 "Invalid request body"
// @Failure 422 {object} ErrorResponse "Invalid exchange request"
// @Failure 404 {object} ErrorResponse "No exchanges found for the given amount and banknotes"
// @Failure 500 "Internal server error"
// @Router /exchange [post]
func (s *exchangeRoutes) exchangeBanknotes(w http.ResponseWriter, r *http.Request) {
	var req ExchangeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.l.Error().Err(err).Msg("Failed to decode request body")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())

		return
	}

	err = validateExchangeRequest(req)
	if err != nil {
		s.l.Debug().Err(err).Msg("Invalid exchange request")
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})

		return
	}

	exchanges, err := s.exchangeService.CalculateExchanges(req.Amount, req.Banknotes)
	if err != nil {
		if errors.Is(err, service.ErrNoExchanges) {
			s.l.Info().Err(err).Msgf("No exchanges found for amount = %d with banknotes %v", req.Amount, req.Banknotes)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
		} else {
			s.l.Error().Err(err).Msg("Failed to calculate exchanges")
			render.Status(r, http.StatusInternalServerError)
		}

		return
	}

	s.l.Info().Msgf("Successful exchange request: %+v", req)
	render.JSON(w, r, ExchangeResponse{Exchanges: exchanges})
}

var (
	ErrRequiredField           = errors.New("field is required")
	ErrBanknotesMustContainOne = errors.New("banknotes field must contain at least one element")
	ErrFieldMinValue           = errors.New("field must be greater than or equal to")
	ErrFieldMustBeUnique       = errors.New("field must contain unique values")
)

func validateExchangeRequest(req ExchangeRequest) error {
	validate := validator.New()

	err := validate.Struct(req)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationErr := range validationErrors {
				switch validationErr.Tag() {
				case "required":
					return fmt.Errorf("%s: %w", validationErr.Field(), ErrRequiredField)
				case "min":
					if validationErr.Field() == "Banknotes" {
						return ErrBanknotesMustContainOne
					}

					return fmt.Errorf("%s: %w %s", validationErr.Field(), ErrFieldMinValue, validationErr.Param())
				case "unique":
					return fmt.Errorf("%s: %w", validationErr.Field(), ErrFieldMustBeUnique)
				default:
					return validationErr
				}
			}
		}

		return fmt.Errorf("failed to validate request: %w", err)
	}

	return nil
}
