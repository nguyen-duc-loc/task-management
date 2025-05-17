package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nguyen-duc-loc/task-management/backend/internal/store"
	"github.com/nguyen-duc-loc/task-management/backend/internal/token"
)

type createTaskRequest struct {
	Name     string `json:"name" binding:"required"`
	Deadline string `json:"deadline" binding:"required,iso8601"`
}

func (s *Server) createTaskHandler(ctx *gin.Context) {
	var req createTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	id, err := gonanoid.New()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	deadline, _ := time.Parse(time.RFC3339, req.Deadline)
	arg := store.CreateTaskParams{
		ID:        id,
		CreatorID: authPayload.UserID,
		Name:      req.Name,
		Deadline:  deadline,
	}

	task, err := s.storage.CreateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}

type getTasksRequest struct {
	Name          string `form:"name" binding:"omitempty"`
	StartDeadline string `form:"start_deadline" binding:"omitempty,iso8601"`
	EndDeadline   string `form:"end_deadline" binding:"omitempty,iso8601"`
	Completed     *bool  `form:"completed" binding:"omitempty"`
	Page          int32  `form:"page" binding:"omitempty,min=1"`
	Limit         int32  `form:"limit" binding:"omitempty,min=1,max=20"`
}

func (s *Server) getTasksHandler(ctx *gin.Context) {
	var req getTasksRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := store.GetTasksParams{
		CreatorID: authPayload.UserID,
		Limit:     req.Limit,
		Offset:    (req.Page - 1) * req.Limit,
	}

	if len(req.Name) > 0 {
		arg.Name = pgtype.Text{
			String: req.Name,
			Valid:  true,
		}
	}

	if len(req.StartDeadline) > 0 {
		startDeadline, _ := time.Parse(time.RFC3339, req.StartDeadline)
		arg.StartDeadline = pgtype.Timestamptz{
			Time:  startDeadline,
			Valid: true,
		}
	}

	if len(req.EndDeadline) > 0 {
		endDeadline, _ := time.Parse(time.RFC3339, req.EndDeadline)
		arg.EndDeadline = pgtype.Timestamptz{
			Time:  endDeadline,
			Valid: true,
		}
	}

	if req.Completed != nil {
		arg.Completed = pgtype.Bool{
			Bool:  *req.Completed,
			Valid: true,
		}
	}

	if arg.Limit == 0 {
		arg.Limit = 5
	}

	if arg.Offset < 0 {
		arg.Offset = 0
	}

	tasks, err := s.storage.GetTasks(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

type getTaskByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (s *Server) getTaskByIDHandler(ctx *gin.Context) {
	var req getTaskByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	task, err := s.storage.GetTaskByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if task.CreatorID != authPayload.UserID {
		err := errors.New("task doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}
