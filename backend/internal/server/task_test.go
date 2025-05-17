package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

type eqGetTasksParamsMatcher struct {
	arg store.GetTasksParams
}

func (e eqGetTasksParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(store.GetTasksParams)
	if !ok {
		return false
	}

	if arg.CreatorID != e.arg.CreatorID {
		return false
	}

	if arg.Name.String != e.arg.Name.String {
		return false
	}

	if arg.Limit != e.arg.Limit {
		return false
	}

	if arg.Offset != e.arg.Offset {
		return false
	}

	if arg.StartDeadline.Time.Sub(e.arg.StartDeadline.Time).Abs() > time.Second {
		return false
	}

	if arg.EndDeadline.Time.Sub(e.arg.EndDeadline.Time).Abs() > time.Second {
		return false
	}

	if arg.Completed.Bool != e.arg.Completed.Bool {
		return false
	}

	return true
}

func (e eqGetTasksParamsMatcher) String() string {
	return fmt.Sprintf("arg: %v", e.arg)
}

func EqGetTasksParams(arg store.GetTasksParams) gomock.Matcher {
	return eqGetTasksParamsMatcher{arg}
}

func requireBodyMatchTasks(t *testing.T, body *bytes.Buffer, tasks []store.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTasks []store.Task
	err = json.Unmarshal(data, &gotTasks)
	require.NoError(t, err)
	require.Equal(t, len(tasks), len(gotTasks))

	for i := range tasks {
		require.Equal(t, tasks[i].ID, gotTasks[i].ID)
		require.Equal(t, tasks[i].Name, gotTasks[i].Name)
		require.Equal(t, tasks[i].CreatorID, gotTasks[i].CreatorID)
		require.Equal(t, tasks[i].Completed, gotTasks[i].Completed)
		require.WithinDuration(t, tasks[i].Deadline, gotTasks[i].Deadline, time.Second)
		require.WithinDuration(t, tasks[i].CreatedAt, gotTasks[i].CreatedAt, time.Second)
	}
}

func TestGetTasksHandler(t *testing.T) {
	user, _ := randomUser(t)

	n := 10
	tasks := make([]store.Task, n)
	for i := 0; i < n; i++ {
		tasks[i] = randomTask(t, user.ID)
	}

	type Query struct {
		Name          string
		StartDeadline string
		EndDeadline   string
		Completed     *bool
		Page          *int32
		Limit         *int32
	}

	incomplete, limit, page := new(bool), new(int32), new(int32)
	*incomplete = false
	*limit = int32(n)
	*page = 1

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				Page:          page,
				Limit:         limit,
				Name:          "",
				StartDeadline: time.Now().Add(-time.Hour).Format(time.RFC3339),
				EndDeadline:   time.Now().Add(time.Hour).Format(time.RFC3339),
				Completed:     incomplete,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.GetTasksParams{
					CreatorID: user.ID,
					Limit:     int32(n),
					Offset:    0,
					Name: pgtype.Text{
						String: "",
						Valid:  true,
					},
					StartDeadline: pgtype.Timestamptz{
						Time:  time.Now().Add(-time.Hour),
						Valid: true,
					},
					EndDeadline: pgtype.Timestamptz{
						Time:  time.Now().Add(+time.Hour),
						Valid: true,
					},
					Completed: pgtype.Bool{
						Bool:  false,
						Valid: true,
					},
				}
				storage.EXPECT().
					GetTasks(gomock.Any(), EqGetTasksParams(arg)).
					Times(1).
					Return(tasks, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTasks(t, recorder.Body, tasks)
			},
		},
		{
			name:  "NoQueryString",
			query: Query{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				arg := store.GetTasksParams{
					CreatorID: user.ID,
					Limit:     5,
					Offset:    0,
				}
				storage.EXPECT().
					GetTasks(gomock.Any(), EqGetTasksParams(arg)).
					Times(1).
					Return(tasks, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTasks(t, recorder.Body, tasks)
			},
		},
		{
			name:  "UnauthorizedUser",
			query: Query{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTasks(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "InternalError",
			query: Query{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTasks(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]store.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidDeadline",
			query: Query{
				StartDeadline: "invalid",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTasks(gomock.Any(), gomock.Any()).
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

			url := "/tasks"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()

			if len(tc.query.Name) > 0 {
				q.Add("name", tc.query.Name)
			}

			if len(tc.query.StartDeadline) > 0 {
				q.Add("start_deadline", tc.query.StartDeadline)
			}

			if len(tc.query.EndDeadline) > 0 {
				q.Add("end_deadline", tc.query.EndDeadline)
			}

			if tc.query.Completed != nil {
				q.Add("completed", strconv.FormatBool(*tc.query.Completed))
			}

			if tc.query.Page != nil {
				q.Add("page", strconv.FormatInt(int64(*tc.query.Page), 10))
			}

			if tc.query.Limit != nil {
				q.Add("limit", strconv.FormatInt(int64(*tc.query.Limit), 10))
			}

			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetTaskByIDHandler(t *testing.T) {
	user, _ := randomUser(t)
	task := randomTask(t, user.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name: "NotFound",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(store.Task{}, store.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 100, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(store.Task{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Eq(task.ID)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				storage.EXPECT().
					GetTaskByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(store.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/tasks/%s", task.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
