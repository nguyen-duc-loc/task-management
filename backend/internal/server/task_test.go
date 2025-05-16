package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	mockdb "github.com/nguyen-duc-loc/task-management/backend/internal/database/mock"
	"github.com/nguyen-duc-loc/task-management/backend/internal/store"
	"github.com/nguyen-duc-loc/task-management/backend/internal/token"
	"github.com/nguyen-duc-loc/task-management/backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomTask(t *testing.T, creatorID int64) store.Task {
	id, err := gonanoid.New()
	require.NoError(t, err)

	return store.Task{
		ID:        id,
		CreatorID: creatorID,
		Name:      util.RandomPrintableString(100),
		Deadline:  time.Now().Add(time.Hour),
	}
}

type eqCreateTaskParamsMatcher struct {
	arg store.CreateTaskParams
}

func (e eqCreateTaskParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(store.CreateTaskParams)
	if !ok {
		return false
	}

	if e.arg.Name != arg.Name {
		return false
	}

	if e.arg.CreatorID != arg.CreatorID {
		return false
	}

	if e.arg.Deadline.Sub(arg.Deadline) > time.Second || arg.Deadline.Sub(e.arg.Deadline) > time.Second {
		return false
	}

	if len(arg.ID) == 0 {
		return false
	}

	return true
}

func (e eqCreateTaskParamsMatcher) String() string {
	return fmt.Sprintf("arg: %v", e.arg)
}

func EqCreateTaskParams(arg store.CreateTaskParams) gomock.Matcher {
	return eqCreateTaskParamsMatcher{arg}
}

func requireBodyMatchTask(t *testing.T, body *bytes.Buffer, task store.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTask store.Task
	err = json.Unmarshal(data, &gotTask)

	require.NoError(t, err)
	require.Equal(t, gotTask.ID, task.ID)
	require.Equal(t, gotTask.CreatorID, task.CreatorID)
	require.Equal(t, gotTask.Name, task.Name)
	require.WithinDuration(t, gotTask.Deadline, task.Deadline, time.Second)
	require.False(t, gotTask.Completed)
}

func TestCreateTaskHandler(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(t, user.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":     task.Name,
				"deadline": task.Deadline,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.CreateTaskParams{
					ID:        task.ID,
					CreatorID: task.CreatorID,
					Name:      task.Name,
					Deadline:  task.Deadline,
				}
				storage.EXPECT().
					CreateTask(gomock.Any(), EqCreateTaskParams(arg)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name: "UnauthorizedUser",
			body: gin.H{
				"name":     task.Name,
				"deadline": task.Deadline,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name":     task.Name,
				"deadline": task.Deadline,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidDeadline",
			body: gin.H{
				"name":     task.Name,
				"deadline": "invalidDeadline",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mockdb.NewMockStorage(ctrl)
			tc.buildStubs(storage)

			server, err := NewServer(storage)
			require.NoError(t, err)
			server.RegisterRoutes()
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/tasks"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
