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
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	Deadline    string `json:"deadline" binding:"required,iso8601"`
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
		Title:     req.Title,
		Deadline:  deadline,
	}
	if len(req.Description) > 0 {
		arg.Description = pgtype.Text{
			String: req.Description,
			Valid:  true,
		}
	}

	task, err := s.storage.CreateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, successResponse(task))
}

type getTasksRequest struct {
	Title         string `form:"title" binding:"omitempty"`
	Description   string `form:"description" binding:"omitempty"`
	StartDeadline string `form:"start_deadline" binding:"omitempty,iso8601"`
	EndDeadline   string `form:"end_deadline" binding:"omitempty,iso8601"`
	Completed     *bool  `form:"completed" binding:"omitempty"`
	Page          int32  `form:"page" binding:"omitempty,min=1"`
	Limit         int32  `form:"limit" binding:"omitempty,min=1,max=20"`
}

type GetTaskRow struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description pgtype.Text `json:"description"`
	CreatorID   int64       `json:"creator_id"`
	Deadline    time.Time   `json:"deadline"`
	Completed   bool        `json:"completed"`
	CreatedAt   time.Time   `json:"created_at"`
}

type getTasksResponse struct {
	Total int64        `json:"total"`
	Tasks []GetTaskRow `json:"tasks"`
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

	if len(req.Title) > 0 {
		arg.Title = pgtype.Text{
			String: req.Title,
			Valid:  true,
		}
	}

	if len(req.Description) > 0 {
		arg.Description = pgtype.Text{
			String: req.Description,
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

	var rsp getTasksResponse

	if len(tasks) == 0 {
		rsp = getTasksResponse{
			Total: 0,
			Tasks: []GetTaskRow{},
		}
	} else {
		rsp = getTasksResponse{
			Total: tasks[0].Total,
			Tasks: []GetTaskRow{},
		}
		for _, task := range tasks {
			rsp.Tasks = append(rsp.Tasks, GetTaskRow{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				CreatorID:   task.CreatorID,
				Deadline:    task.Deadline,
				Completed:   task.Completed,
				CreatedAt:   task.CreatedAt,
			})
		}
	}

	ctx.JSON(http.StatusOK, successResponse(rsp))
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

	ctx.JSON(http.StatusOK, successResponse(task))
}

type updateTaskRequest struct {
	ID          string  `uri:"id" binding:"required"`
	Title       string  `json:"title" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
	Deadline    string  `json:"deadline" binding:"omitempty,iso8601"`
	Completed   *bool   `json:"completed" binding:"omitempty"`
}

func (s *Server) updateTasksHandler(ctx *gin.Context) {
	var req updateTaskRequest
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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := store.UpdateTaskParams{
		ID: req.ID,
	}

	if len(req.Title) > 0 {
		arg.Title = pgtype.Text{
			String: req.Title,
			Valid:  true,
		}
	}

	if req.Description != nil {
		arg.Description = pgtype.Text{
			String: *req.Description,
			Valid:  true,
		}
	}

	if len(req.Deadline) > 0 {
		deadline, _ := time.Parse(time.RFC3339, req.Deadline)
		arg.Deadline = pgtype.Timestamptz{
			Time:  deadline,
			Valid: true,
		}
	}

	if req.Completed != nil {
		arg.Completed = pgtype.Bool{
			Bool:  *req.Completed,
			Valid: true,
		}
	}

	newTask, err := s.storage.UpdateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(newTask))
}

type deleteTaskRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (s *Server) deleteTaskHandler(ctx *gin.Context) {
	var req deleteTaskRequest
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

	err = s.storage.DeleteTask(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(nil))
}
