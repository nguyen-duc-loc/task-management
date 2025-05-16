package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
