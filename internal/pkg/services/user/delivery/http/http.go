package http

import (
	"TaskTracker/internal/pkg/constans"
	"TaskTracker/internal/pkg/services/user"
	"TaskTracker/pkg/responser"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

type Params struct {
	fx.In

	Usecase user.Usecase
	Logger  *slog.Logger
}

type Handler struct {
	useCase user.Usecase
	log     *slog.Logger
}

func New(params Params) *Handler {
	return &Handler{
		useCase: params.Usecase,
		log:     params.Logger,
	}
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.Context().Value(constans.ContextValue).(uuid.UUID)

	userData, err := h.useCase.GetUserByID(r.Context(), idStr)
	if err != nil {
		h.log.Error("get user", "error", err.Error())
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, userData)
}
