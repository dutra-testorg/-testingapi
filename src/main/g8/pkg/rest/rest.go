package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Gympass/gcore/v3/gcontext"
	"github.com/Gympass/gcore/v3/gerror"
	"github.com/Gympass/gcore/v3/glog"
	uuid "github.com/gofrs/uuid"
	"github.com/pkg/errors"

	// only for swagger dependency
	_ "github.com/Gympass/gcore/v3/ghandler"
)

// Rest interface
type Rest struct {
	service string
	logger  glog.Logger
}

// Config used by rest package
type Config struct {
	Service string
	Logger  glog.Logger
}

// New creates a rest struct based on configuration properties
func New(cfg Config) *Rest {
	return &Rest{
		service: cfg.Service,
		logger:  cfg.Logger,
	}
}

// SendJSON is a helper function to send a JSON as an HTTP response.
// It sets the header Content-Type as application/json.
// If it fails to write the JSON an internal server error is generated.
func SendJSON(w http.ResponseWriter, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.Wrap(err, "encoding json")
	}

	return nil
}

// DeserializeJSON is a helper function to read a JSON from request body.
// PS.: payload needs to be a pointer.
func DeserializeJSON(r *http.Request, payload interface{}) error {
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		return gerror.NewBadRequest(err).WithMessage("bad json format")
	}

	return nil
}

// @title Swagger No-API
// @version 1.0
// @description This is a health-check generated with 'swag init -g ./internal/rest/rest.go'.
// @termsOfService http://swagger.io/terms/
// @contact.name Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// Health ...
// ShowEntity godoc
// @Summary Provide health-check endpoint
// @Description Health-check returning a dummy HealthCheckResponse (config)
// @Produce  json
// @Success 200 {object} HealthCheckResponse
// @Failure 500 {object} ghandler.HTTPError
// @Router /health [get]
func (c Rest) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u4, err := uuid.NewV4()
	if err != nil {
		gcontext.AddError(ctx, err)
		c.logger.Error(ctx, "Failed to generate UUID.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = SendJSON(
		w,
		HealthCheckResponse{
			ID:      u4,
			Service: c.service,
		},
	)
	if err != nil {
		gcontext.AddError(ctx, err)
		c.logger.Error(ctx, "Failed to write JSON response.")
	}
}

// HealthCheckResponse ...
type HealthCheckResponse struct {
	ID      uuid.UUID `json:"id"`
	Service string    `json:"service"`
}
