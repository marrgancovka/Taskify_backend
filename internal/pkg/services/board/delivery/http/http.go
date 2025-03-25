package http

import (
	"TaskTracker/internal/models"
	"TaskTracker/internal/pkg/constans"
	"TaskTracker/internal/pkg/services/board"
	"TaskTracker/pkg/reader"
	"TaskTracker/pkg/responser"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

type Params struct {
	fx.In

	Usecase board.Usecase
	Logger  *slog.Logger
}

type Handler struct {
	useCase board.Usecase
	log     *slog.Logger
}

func New(params Params) *Handler {
	return &Handler{
		useCase: params.Usecase,
		log:     params.Logger,
	}
}

func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constans.ContextValue).(uuid.UUID)

	newBoard := &models.Board{}
	if err := reader.ReadRequestData(r, newBoard); err != nil {
		responser.Send400(w, "некорректные данные")
		return
	}
	newBoard.OwnerID = userId

	newBoard, err := h.useCase.CreateBoard(r.Context(), newBoard)
	if err != nil {
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, newBoard)
}

func (h *Handler) GetUserListBoard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constans.ContextValue).(uuid.UUID)

	boardList, err := h.useCase.GetUserListBoards(r.Context(), userId)
	if err != nil {
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, boardList)
}
