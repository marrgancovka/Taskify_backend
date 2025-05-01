package http

import (
	"TaskTracker/internal/models"
	"TaskTracker/internal/pkg/constans"
	"TaskTracker/internal/pkg/services/board"
	"TaskTracker/pkg/reader"
	"TaskTracker/pkg/responser"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	newBoardData := &models.Board{}
	if err := reader.ReadRequestData(r, newBoardData); err != nil {
		responser.Send400(w, "некорректные данные")
		return
	}
	newBoardData.OwnerID = userId

	newBoard, err := h.useCase.CreateBoard(r.Context(), newBoardData)
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

func (h *Handler) SetFavouriteBoard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constans.ContextValue).(uuid.UUID)
	vars := mux.Vars(r)
	boardIDStr := vars["boardID"]
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}
	err = h.useCase.SetFavouriteBoard(r.Context(), boardID, userId)
	if err != nil {
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, nil)

}

func (h *Handler) SetNoFavouriteBoard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constans.ContextValue).(uuid.UUID)
	vars := mux.Vars(r)
	boardIDStr := vars["boardID"]
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}
	err = h.useCase.SetNoFavouriteBoard(r.Context(), boardID, userId)
	if err != nil {
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, nil)
	return
}

func (h *Handler) GetBoardTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constans.ContextValue).(uuid.UUID)
	vars := mux.Vars(r)
	boardIDStr := vars["boardID"]
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
	}
	task, err := h.useCase.GetTaskInBoard(r.Context(), boardID, userId)
	if err != nil {
		h.log.Error("getTaskInBoard", "err", err.Error())
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, task)
}

//func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
//	boardID, err := reader.ReadVarsUUID(r, "boardID")
//	if err != nil {
//		http.Error(w, "Invalid board ID", http.StatusBadRequest)
//		return
//	}
//	memberData := &models.BoardMemberAdd{}
//	if err = reader.ReadRequestData(r, &memberData); err != nil {
//		http.Error(w, "Invalid board ID", http.StatusBadRequest)
//		return
//	}
//	memberData.BoardID = boardID
//	err = h.useCase.AddMember(r.Context(), memberData)
//	if err != nil {
//		h.log.Error("addMember", "err", err.Error())
//		responser.Send500(w, err.Error())
//		return
//	}
//	responser.Send200(w, nil)
//}

func (h *Handler) CreateSection(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constans.ContextValue).(uuid.UUID)

	sectionData := &models.Section{}
	if err := reader.ReadRequestData(r, sectionData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.log.Debug("CreateSection", "userId", userId, "sectionData", sectionData)

	section, err := h.useCase.CreateSection(r.Context(), sectionData, userId)
	if err != nil {
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, section)

}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	taskData := &models.TaskCreate{}
	if err := reader.ReadRequestData(r, taskData); err != nil {
		h.log.Error("CreateTask", "err", err.Error())
		responser.Send400(w, err.Error())
		return
	}
	userID := r.Context().Value(constans.ContextValue).(uuid.UUID)
	createdTask, err := h.useCase.AddTask(r.Context(), taskData, userID)
	if err != nil {
		h.log.Error("addTask", "err", err.Error())
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, &createdTask)
}

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constans.ContextValue).(uuid.UUID)
	boardID, err := reader.ReadVarsUUID(r, "boardID")
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	addMemberData := &models.BoardMemberAdd{}
	if err := reader.ReadRequestData(r, addMemberData); err != nil {
		h.log.Error("AddMember", "err", err.Error())
		responser.Send400(w, err.Error())
		return
	}
	addMemberData.BoardID = boardID

	addedMember, err := h.useCase.AddMember(r.Context(), addMemberData, userId)
	if err != nil {
		h.log.Error("addMember", "err", err.Error())
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, addedMember)

}

func (h *Handler) GetBoardMembers(w http.ResponseWriter, r *http.Request) {
	boardID, err := reader.ReadVarsUUID(r, "boardID")
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	list, err := h.useCase.GetBoardMemberList(r.Context(), boardID)
	if err != nil {
		h.log.Error("GetBoardMembers", "err", err.Error())
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, list)
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	boardID, err := reader.ReadVarsUUID(r, "boardID")
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	tasks, err := h.useCase.GetAllTasks(r.Context(), boardID)
	if err != nil {
		h.log.Error("GetAllTasks", "err", err.Error())
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, tasks)

}

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	updateData := &models.UpdateTask{}
	if err := reader.ReadRequestData(r, updateData); err != nil {
		h.log.Error("EditTask", "err", err.Error())
		responser.Send400(w, err.Error())
		return
	}
	userId := r.Context().Value(constans.ContextValue).(uuid.UUID)

	taskID, err := h.useCase.UpdateTask(r.Context(), updateData, userId)
	if err != nil {
		h.log.Error("EditTask", "err", err.Error())
		responser.Send500(w, err.Error())
		return
	}
	responser.Send200(w, taskID)
}
