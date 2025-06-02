package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nguyen-duc-loc/task-management/backend/internal/store"
	"github.com/nguyen-duc-loc/task-management/backend/util"
)

var (
	errUsernameConflict   = errors.New("username already exists")
	errInvalidCredentials = errors.New("invalid credentials")
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user store.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

func (s *Server) createUserHandler(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := store.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
	}

	user, err := s.storage.CreateUser(ctx, arg)
	if err != nil {
		if store.ErrorCode(err) == store.UniqueViolation {
			ctx.JSON(http.StatusConflict, errorResponse(errUsernameConflict))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, successResponse(newUserResponse(user)))
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	AccessToken         string       `json:"access_token"`
	AccessTokenExpireAt time.Time    `json:"access_token_expire_at"`
	User                userResponse `json:"user"`
}

func (s *Server) loginUserHandler(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.storage.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errInvalidCredentials))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errInvalidCredentials))
		return
	}

	jwtConfig, err := util.LoadJWTConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		jwtConfig.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(loginResponse{
		AccessToken:         accessToken,
		AccessTokenExpireAt: accessPayload.ExpireAt,
		User:                newUserResponse(user),
	}))
}
