package micro

import (
	"context"
	"net/http"

	"github.com/gympass/$name;format="lower,hyphen"$/pkg/rest"

	"github.com/Gympass/gcore/v3/gerror"
	"github.com/Gympass/gcore/v3/glog"

	// only for swagger dependency
	_ "github.com/Gympass/gcore/v3/ghandler"
)

// Handler struct
type Handler struct {
	svc    *Service
	logger glog.Logger
}

// NewHandler holding base struct
func NewHandler(s *Service, l glog.Logger) *Handler {
	return &Handler{svc: s, logger: l}
}

// Demo ...
// ShowEntity godoc
// @Summary Get demo
// @Description demo endpoint returning a Demo struct
// @Param uid path string true "uuidv4 (UUIDv4)"
// @Produce  json
// @Success 200 {object} Demo
// @Failure 400 {object} ghandler.HTTPError
// @Failure 500 {object} ghandler.HTTPError
// @Router /v1/demo/{uid} [get]
func (h *Handler) Demo(w http.ResponseWriter, r *http.Request) error {
	uid, err := rest.GetUUID(r, "uid")
	if err != nil {
		return gerror.NewBadRequest(ErrInvalidID).WithMessage(ErrInvalidID.Error())
	}

	demo, err := h.svc.Demo(context.Background(), uid.String())
	if err != nil {
		return gerror.NewInternalServerError(ErrInternal).WithMessage(ErrInternal.Error())
	}

	return rest.SendJSON(w, &demo)
}
